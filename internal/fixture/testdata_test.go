package fixture

import (
	"os"
	"path/filepath"
	"testing"
)

func minimalFixture() FixtureSet {
	return FixtureSet{
		SchemaVersion: FixtureSchemaVersion,
		ID:            "fixture-minimal",
		Scenes: []Scene{
			{
				ID:     "terminal-shell",
				Family: FamilyTerminalAgent,
				Regions: []Region{
					{
						ID:         "shell-output",
						Kind:       RegionSurface,
						State:      StateActive,
						Background: RoleTerminalBackground,
						Foreground: RoleTerminalForeground,
						Blocks: []ContentBlock{
							{
								ID:         "command-result",
								Kind:       ContentCode,
								State:      StateSuccess,
								Text:       "ok",
								Foreground: RoleSuccess,
							},
						},
					},
				},
			},
		},
	}
}

func minimalPalette(id PaletteID) Palette {
	return Palette{
		SchemaVersion: PaletteSchemaVersion,
		ID:            id,
		SemanticCore:  authoredRoles(semanticCoreRoleIDs),
		Terminal: TerminalPalette{
			ANSI: authoredRoles(ansiRoleIDs),
			Aliases: []Role{
				aliasRole(RoleTerminalBackground, RoleBackground),
				aliasRole(RoleTerminalForeground, RoleForeground),
				aliasRole(RoleTerminalCursor, RoleCursor),
				aliasRole(RoleTerminalSelectionBackground, RoleSelectionBackground),
				aliasRole(RoleTerminalSelectionForeground, RoleSelectionForeground),
			},
		},
	}
}

func authoredRoles(ids []RoleID) []Role {
	roles := make([]Role, 0, len(ids))
	for _, id := range ids {
		roles = append(roles, Role{
			ID: id,
			Value: &AuthoredColor{
				OKLCH: OKLCH{Lightness: 0.5, Chroma: 0, Hue: 0},
				SRGB:  "#777777",
				Provenance: Provenance{
					ProfileClaims: []ProfileClaimID{"claim-synthetic"},
					Judgments:     []JudgmentID{"judgment-synthetic"},
				},
				Validity: ValidityMetadata{
					Evidence: []ValidityEvidenceID{"validity-synthetic"},
				},
			},
		})
	}
	return roles
}

func aliasRole(id, target RoleID) Role {
	return Role{ID: id, Alias: &Alias{Target: target}}
}

func readGolden(t *testing.T, name string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read golden %s: %v", name, err)
	}
	return string(contents)
}
