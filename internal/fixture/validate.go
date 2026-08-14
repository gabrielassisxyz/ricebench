package fixture

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

var (
	stableIDPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)
	hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
)

func (definition FixtureSet) Validate() error {
	if definition.SchemaVersion != FixtureSchemaVersion {
		return fmt.Errorf("schemaVersion: got %q, want %q", definition.SchemaVersion, FixtureSchemaVersion)
	}
	if err := validateID("id", string(definition.ID)); err != nil {
		return err
	}
	if len(definition.Scenes) == 0 {
		return fmt.Errorf("scenes: must contain at least one scene")
	}

	sceneIDs := make(map[SceneID]struct{}, len(definition.Scenes))
	regionIDs := make(map[RegionID]struct{})
	blockIDs := make(map[ContentBlockID]struct{})
	for sceneIndex, scene := range definition.Scenes {
		path := fmt.Sprintf("scenes[%d]", sceneIndex)
		if err := validateUniqueID(path+".id", scene.ID, sceneIDs); err != nil {
			return err
		}
		if !validSceneFamily(scene.Family) {
			return fmt.Errorf("%s.family: unknown %q", path, scene.Family)
		}
		if len(scene.Regions) == 0 {
			return fmt.Errorf("%s.regions: must contain at least one region", path)
		}
		for regionIndex, region := range scene.Regions {
			regionPath := fmt.Sprintf("%s.regions[%d]", path, regionIndex)
			if err := validateUniqueID(regionPath+".id", region.ID, regionIDs); err != nil {
				return err
			}
			if !validRegionKind(region.Kind) {
				return fmt.Errorf("%s.kind: unknown %q", regionPath, region.Kind)
			}
			if !validSemanticState(region.State) {
				return fmt.Errorf("%s.state: unknown %q", regionPath, region.State)
			}
			if err := validateRoleReferences(regionPath, region.Background, region.Foreground, region.Border); err != nil {
				return err
			}
			for blockIndex, block := range region.Blocks {
				blockPath := fmt.Sprintf("%s.blocks[%d]", regionPath, blockIndex)
				if err := validateUniqueID(blockPath+".id", block.ID, blockIDs); err != nil {
					return err
				}
				if !validContentKind(block.Kind) {
					return fmt.Errorf("%s.kind: unknown %q", blockPath, block.Kind)
				}
				if !validSemanticState(block.State) {
					return fmt.Errorf("%s.state: unknown %q", blockPath, block.State)
				}
				if block.Text != "" && len(block.Runs) != 0 {
					return fmt.Errorf("%s: text and runs are mutually exclusive", blockPath)
				}
				if err := validateRoleReferences(blockPath, block.Background, block.Foreground, block.Border); err != nil {
					return err
				}
				for runIndex, run := range block.Runs {
					runPath := fmt.Sprintf("%s.runs[%d]", blockPath, runIndex)
					if run.Text == "" {
						return fmt.Errorf("%s.text: must not be empty", runPath)
					}
					if !validSemanticState(run.State) {
						return fmt.Errorf("%s.state: unknown %q", runPath, run.State)
					}
					if err := validateRoleReferences(runPath, run.Background, run.Foreground); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (palette Palette) Validate() error {
	if palette.SchemaVersion != PaletteSchemaVersion {
		return fmt.Errorf("schemaVersion: got %q, want %q", palette.SchemaVersion, PaletteSchemaVersion)
	}
	if err := validateID("id", string(palette.ID)); err != nil {
		return err
	}

	allRoles := make(map[RoleID]Role, len(semanticCoreRoleIDs)+len(ansiRoleIDs)+len(terminalAliasRoleIDs))
	if err := validateRoleGroup("semanticCore", palette.SemanticCore, semanticCoreRoleIDs, true, allRoles); err != nil {
		return err
	}
	if err := validateRoleGroup("terminal.ansi", palette.Terminal.ANSI, ansiRoleIDs, true, allRoles); err != nil {
		return err
	}
	if err := validateRoleGroup("terminal.aliases", palette.Terminal.Aliases, terminalAliasRoleIDs, false, allRoles); err != nil {
		return err
	}

	semantic := roleSet(semanticCoreRoleIDs)
	for index, role := range palette.Terminal.Aliases {
		resolved, err := resolveRole(role.ID, allRoles, nil)
		if err != nil {
			return fmt.Errorf("terminal.aliases[%d].alias.target: %w", index, err)
		}
		if _, ok := semantic[resolved]; !ok {
			return fmt.Errorf("terminal alias %q resolves to non-semantic role %q", role.ID, resolved)
		}
	}
	return nil
}

func ValidatePair(definition FixtureSet, palette Palette) error {
	if err := definition.Validate(); err != nil {
		return fmt.Errorf("fixture %q: %w", definition.ID, err)
	}
	if err := palette.Validate(); err != nil {
		return fmt.Errorf("palette %q: %w", palette.ID, err)
	}

	roles := paletteRoleMap(palette)
	for sceneIndex, scene := range definition.Scenes {
		for regionIndex, region := range scene.Regions {
			path := fmt.Sprintf("scenes[%d].regions[%d]", sceneIndex, regionIndex)
			if err := ensureRolesExist(path, palette.ID, roles, region.Background, region.Foreground, region.Border); err != nil {
				return err
			}
			for blockIndex, block := range region.Blocks {
				blockPath := fmt.Sprintf("%s.blocks[%d]", path, blockIndex)
				if err := ensureRolesExist(blockPath, palette.ID, roles, block.Background, block.Foreground, block.Border); err != nil {
					return err
				}
				for runIndex, run := range block.Runs {
					runPath := fmt.Sprintf("%s.runs[%d]", blockPath, runIndex)
					if err := ensureRolesExist(runPath, palette.ID, roles, run.Background, run.Foreground); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func validateRoleGroup(path string, roles []Role, required []RoleID, authored bool, all map[RoleID]Role) error {
	for index, role := range roles {
		rolePath := fmt.Sprintf("%s[%d]", path, index)
		if err := validateID(rolePath+".id", string(role.ID)); err != nil {
			return err
		}
		if _, exists := all[role.ID]; exists {
			return fmt.Errorf("%s.id: duplicate %q", rolePath, role.ID)
		}
		if (role.Value == nil) == (role.Alias == nil) {
			return fmt.Errorf("%s: role %q must contain exactly one of value or alias", rolePath, role.ID)
		}
		if authored && role.Value == nil {
			return fmt.Errorf("%s: role %q must be authored", rolePath, role.ID)
		}
		if !authored && role.Alias == nil {
			return fmt.Errorf("%s: role %q must be an alias", rolePath, role.ID)
		}
		if role.Value != nil {
			if err := validateAuthoredColor(rolePath+".value", *role.Value); err != nil {
				return err
			}
		}
		all[role.ID] = role
	}

	present := make(map[RoleID]struct{}, len(roles))
	for _, role := range roles {
		present[role.ID] = struct{}{}
	}
	for _, requiredRole := range required {
		if _, ok := present[requiredRole]; !ok {
			return fmt.Errorf("%s: missing required role %q", path, requiredRole)
		}
	}
	for _, role := range roles {
		if !containsRole(required, role.ID) {
			return fmt.Errorf("%s: role %q is not part of schema version 1", path, role.ID)
		}
	}
	return nil
}

func validateAuthoredColor(path string, color AuthoredColor) error {
	// This layer validates the wire contract only. The authoritative bounded types,
	// gamut mapping, conversion and OKLCH-to-sRGB consistency belong to the color package.
	if !finite(color.OKLCH.Lightness) || !finite(color.OKLCH.Chroma) || !finite(color.OKLCH.Hue) {
		return fmt.Errorf("%s.oklch: values must be finite", path)
	}
	if !hexColorPattern.MatchString(color.SRGB) {
		return fmt.Errorf("%s.srgb: %q is not a six-digit hexadecimal color", path, color.SRGB)
	}
	if err := validateIDList(path+".provenance.profileClaims", color.Provenance.ProfileClaims); err != nil {
		return err
	}
	if err := validateIDList(path+".provenance.judgments", color.Provenance.Judgments); err != nil {
		return err
	}
	return validateIDList(path+".validity.evidence", color.Validity.Evidence)
}

func resolveRole(id RoleID, roles map[RoleID]Role, path []RoleID) (RoleID, error) {
	role, exists := roles[id]
	if !exists {
		return "", fmt.Errorf("role %q does not exist", id)
	}
	for index, visited := range path {
		if visited == id {
			cycle := append(append([]RoleID{}, path[index:]...), id)
			parts := make([]string, len(cycle))
			for partIndex, roleID := range cycle {
				parts[partIndex] = string(roleID)
			}
			return "", fmt.Errorf("alias cycle: %s", strings.Join(parts, " -> "))
		}
	}
	if role.Value != nil {
		return id, nil
	}
	return resolveRole(role.Alias.Target, roles, append(path, id))
}

func paletteRoleMap(palette Palette) map[RoleID]Role {
	roles := make(map[RoleID]Role, len(palette.SemanticCore)+len(palette.Terminal.ANSI)+len(palette.Terminal.Aliases))
	for _, group := range [][]Role{palette.SemanticCore, palette.Terminal.ANSI, palette.Terminal.Aliases} {
		for _, role := range group {
			roles[role.ID] = role
		}
	}
	return roles
}

func ensureRolesExist(path string, paletteID PaletteID, roles map[RoleID]Role, references ...RoleID) error {
	for _, reference := range references {
		if reference == "" {
			continue
		}
		if _, exists := roles[reference]; !exists {
			return fmt.Errorf("%s: role %q does not exist in palette %q", path, reference, paletteID)
		}
	}
	return nil
}

func validateRoleReferences(path string, references ...RoleID) error {
	for _, reference := range references {
		if reference == "" {
			continue
		}
		if err := validateID(path+".role", string(reference)); err != nil {
			return err
		}
	}
	return nil
}

func validateUniqueID[T ~string](path string, id T, seen map[T]struct{}) error {
	if err := validateID(path, string(id)); err != nil {
		return err
	}
	if _, exists := seen[id]; exists {
		return fmt.Errorf("%s: duplicate %q", path, id)
	}
	seen[id] = struct{}{}
	return nil
}

func validateIDList[T ~string](path string, ids []T) error {
	seen := make(map[T]struct{}, len(ids))
	for index, id := range ids {
		if err := validateUniqueID(fmt.Sprintf("%s[%d]", path, index), id, seen); err != nil {
			return err
		}
	}
	return nil
}

func validateID(path, id string) error {
	if !stableIDPattern.MatchString(id) {
		return fmt.Errorf("%s: %q is not a stable ID", path, id)
	}
	return nil
}

func validSceneFamily(value SceneFamily) bool {
	return value == FamilyTerminalAgent || value == FamilyCodeDiff ||
		value == FamilyDesktopShell || value == FamilyReadingMonitor
}

func validRegionKind(value RegionKind) bool {
	switch value {
	case RegionFrame, RegionSurface, RegionRow, RegionColumn, RegionTabs,
		RegionList, RegionTable, RegionStatus, RegionOverlay:
		return true
	default:
		return false
	}
}

func validContentKind(value ContentKind) bool {
	switch value {
	case ContentText, ContentCode, ContentTab, ContentListItem, ContentTableCell,
		ContentStatusItem, ContentSelection, ContentFocusRing, ContentCursor, ContentBadge:
		return true
	default:
		return false
	}
}

func validSemanticState(value SemanticState) bool {
	switch value {
	case StateDefault, StateActive, StateInactive, StateFocused, StateSelected,
		StateMuted, StateSuccess, StateWarning, StateError, StateInfo, StateUrgent,
		StateAdded, StateRemoved, StateModified, StateSearchMatch:
		return true
	default:
		return false
	}
}

func roleSet(ids []RoleID) map[RoleID]struct{} {
	set := make(map[RoleID]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return set
}

func containsRole(roles []RoleID, wanted RoleID) bool {
	for _, role := range roles {
		if role == wanted {
			return true
		}
	}
	return false
}

func finite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
