package fixture

import "fmt"

// ReferenceLevel names the structural level at which a role reference appears.
type ReferenceLevel string

const (
	LevelRegion ReferenceLevel = "region"
	LevelBlock  ReferenceLevel = "block"
	LevelRun    ReferenceLevel = "run"
)

// RoleReference records one role reference in the fixture structure, with its full
// canonical path and the level at which it appears.
type RoleReference struct {
	Role  RoleID
	Path  string
	Level ReferenceLevel
}

// Coverage returns every role reference in the fixture set, with its full path. The
// traversal covers Region, ContentBlock, and ContentRun references together.
func Coverage(definition FixtureSet) []RoleReference {
	var references []RoleReference
	for sceneIndex, scene := range definition.Scenes {
		for regionIndex, region := range scene.Regions {
			path := fmt.Sprintf("scenes[%d].regions[%d]", sceneIndex, regionIndex)
			references = appendRoleReference(references, path+".background", region.Background, LevelRegion)
			references = appendRoleReference(references, path+".foreground", region.Foreground, LevelRegion)
			references = appendRoleReference(references, path+".border", region.Border, LevelRegion)
			for blockIndex, block := range region.Blocks {
				blockPath := fmt.Sprintf("%s.blocks[%d]", path, blockIndex)
				references = appendRoleReference(references, blockPath+".background", block.Background, LevelBlock)
				references = appendRoleReference(references, blockPath+".foreground", block.Foreground, LevelBlock)
				references = appendRoleReference(references, blockPath+".border", block.Border, LevelBlock)
				for runIndex, run := range block.Runs {
					runPath := fmt.Sprintf("%s.runs[%d]", blockPath, runIndex)
					references = appendRoleReference(references, runPath+".background", run.Background, LevelRun)
					references = appendRoleReference(references, runPath+".foreground", run.Foreground, LevelRun)
				}
			}
		}
	}
	return references
}

func appendRoleReference(references []RoleReference, path string, role RoleID, level ReferenceLevel) []RoleReference {
	if role == "" {
		return references
	}
	return append(references, RoleReference{Role: role, Path: path, Level: level})
}

// terminalAliasTargets maps each terminal alias role to the semantic role it names.
// The terminal family references the alias (terminal-cursor) rather than the semantic
// role (cursor), so coverage resolves the alias to prove the semantic role is exercised.
var terminalAliasTargets = map[RoleID]RoleID{
	RoleTerminalBackground:          RoleBackground,
	RoleTerminalForeground:          RoleForeground,
	RoleTerminalCursor:              RoleCursor,
	RoleTerminalSelectionBackground: RoleSelectionBackground,
	RoleTerminalSelectionForeground: RoleSelectionForeground,
}

// RequiredRoleCoverage verifies every required role (semantic core, ANSI, and terminal
// alias) is referenced somewhere across the given sets. A terminal alias reference
// covers its semantic target, so terminal-cursor covers cursor.
func RequiredRoleCoverage(sets ...FixtureSet) error {
	covered := make(map[RoleID]struct{})
	for _, set := range sets {
		for _, reference := range Coverage(set) {
			covered[reference.Role] = struct{}{}
			if target, ok := terminalAliasTargets[reference.Role]; ok {
				covered[target] = struct{}{}
			}
		}
	}
	for _, role := range requiredCoverageRoles() {
		if _, ok := covered[role]; !ok {
			return fmt.Errorf("required role %q is not covered by any Region, ContentBlock, or ContentRun reference", role)
		}
	}
	return nil
}

// RoleCoverage groups every role reference by role ID, resolving terminal aliases to
// their semantic target so a role covered only through an alias still appears under its
// semantic name. Roles with no references are absent from the map.
func RoleCoverage(definition FixtureSet) map[RoleID][]RoleReference {
	byRole := make(map[RoleID][]RoleReference)
	for _, reference := range Coverage(definition) {
		byRole[reference.Role] = append(byRole[reference.Role], reference)
		if target, ok := terminalAliasTargets[reference.Role]; ok {
			byRole[target] = append(byRole[target], reference)
		}
	}
	return byRole
}

// RoleCoverageEntry pairs a required role with the references that cover it.
type RoleCoverageEntry struct {
	Role       RoleID
	References []RoleReference
}

// CoverageView returns every required role in schema order, each with the references
// that cover it across the given sets. A role with no references appears with an empty
// list, so the view shows the full required surface rather than only what is covered.
func CoverageView(sets ...FixtureSet) []RoleCoverageEntry {
	byRole := make(map[RoleID][]RoleReference)
	for _, set := range sets {
		for role, references := range RoleCoverage(set) {
			byRole[role] = append(byRole[role], references...)
		}
	}
	entries := make([]RoleCoverageEntry, 0, len(requiredCoverageRoles()))
	for _, role := range requiredCoverageRoles() {
		entries = append(entries, RoleCoverageEntry{Role: role, References: byRole[role]})
	}
	return entries
}

func requiredCoverageRoles() []RoleID {
	roles := make([]RoleID, 0, len(semanticCoreRoleIDs)+len(ansiRoleIDs)+len(terminalAliasRoleIDs))
	roles = append(roles, semanticCoreRoleIDs...)
	roles = append(roles, ansiRoleIDs...)
	roles = append(roles, terminalAliasRoleIDs...)
	return roles
}
