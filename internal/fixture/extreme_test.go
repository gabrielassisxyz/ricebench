package fixture

import (
	"strings"
	"testing"
)

func TestExtremePalettesValidate(t *testing.T) {
	for _, extreme := range ExtremePalettes() {
		t.Run(string(extreme.ID), func(t *testing.T) {
			if err := extreme.Palette.Validate(); err != nil {
				t.Fatalf("Validate: %v", err)
			}
		})
	}
}

func TestExtremePalettesValidateAgainstEveryFamily(t *testing.T) {
	for _, extreme := range ExtremePalettes() {
		for _, set := range allFixtureSets() {
			t.Run(string(extreme.ID)+"/"+string(set.ID), func(t *testing.T) {
				if err := ValidatePair(set, extreme.Palette); err != nil {
					t.Fatalf("ValidatePair: %v", err)
				}
			})
		}
	}
}

func TestExtremePalettesExposeIntendedFailure(t *testing.T) {
	tests := []struct {
		id       PaletteID
		collapse func(map[RoleID]string) bool
	}{
		{
			id: "extreme-flat-hierarchy",
			collapse: func(c map[RoleID]string) bool {
				return c[RoleBackground] == c[RoleSurface] && c[RoleSurface] == c[RoleSurfaceElevated] &&
					c[RoleForeground] == c[RoleTextSecondary] && c[RoleTextSecondary] == c[RoleTextMuted]
			},
		},
		{
			id: "extreme-collapsed-surfaces",
			collapse: func(c map[RoleID]string) bool {
				return c[RoleSurface] == c[RoleSurfaceElevated]
			},
		},
		{
			id: "extreme-invisible-cursor-selection",
			collapse: func(c map[RoleID]string) bool {
				return c[RoleCursor] == c[RoleBackground] && c[RoleSelectionBackground] == c[RoleBackground]
			},
		},
		{
			id: "extreme-semantic-collision",
			collapse: func(c map[RoleID]string) bool {
				return c[RoleError] == c[RoleSuccess] && c[RoleWarning] == c[RoleInfo]
			},
		},
		{
			id: "extreme-ansi-bright-collapse",
			collapse: func(c map[RoleID]string) bool {
				for index := 0; index < 8; index++ {
					if c[ansiRoleIDs[index]] != c[ansiRoleIDs[index+8]] {
						return false
					}
				}
				return true
			},
		},
	}

	for _, test := range tests {
		t.Run(string(test.id), func(t *testing.T) {
			colors, err := ResolvedColors(extremePaletteByID(t, test.id))
			if err != nil {
				t.Fatalf("ResolvedColors: %v", err)
			}
			if !test.collapse(colors) {
				t.Fatalf("palette %q does not expose its intended failure", test.id)
			}
		})
	}
}

func TestExtremePaletteMissingRoleRejected(t *testing.T) {
	palette := extremePaletteByID(t, "extreme-flat-hierarchy")
	palette.SemanticCore = removeRole(palette.SemanticCore, RoleBackground)

	err := palette.Validate()
	if err == nil {
		t.Fatal("palette missing a required role was accepted")
	}
	if !strings.Contains(err.Error(), string(RoleBackground)) {
		t.Fatalf("error %q does not name the missing role %q", err, RoleBackground)
	}
}

func extremePaletteByID(t *testing.T, id PaletteID) Palette {
	t.Helper()
	for _, extreme := range ExtremePalettes() {
		if extreme.ID == id {
			return extreme.Palette
		}
	}
	t.Fatalf("extreme palette %q not found", id)
	return Palette{}
}

func removeRole(roles []Role, id RoleID) []Role {
	kept := make([]Role, 0, len(roles))
	for _, role := range roles {
		if role.ID != id {
			kept = append(kept, role)
		}
	}
	return kept
}
