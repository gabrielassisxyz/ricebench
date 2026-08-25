package fixture

import "testing"

func TestResolvedColorsFollowsAliases(t *testing.T) {
	palette := minimalPalette("palette-one")
	colors, err := ResolvedColors(palette)
	if err != nil {
		t.Fatalf("ResolvedColors: %v", err)
	}

	// minimalPalette authors every role as #777777, so the alias resolves to the same
	// value as its semantic target rather than to a distinct alias color.
	if colors[RoleTerminalBackground] != colors[RoleBackground] {
		t.Fatalf("terminal-background = %q, want background value %q", colors[RoleTerminalBackground], colors[RoleBackground])
	}
	if colors[RoleBackground] != "#777777" {
		t.Fatalf("background = %q, want #777777", colors[RoleBackground])
	}
}

func TestResolvedColorsRejectsAliasToMissingRole(t *testing.T) {
	palette := minimalPalette("palette-one")
	palette.Terminal.Aliases = append(palette.Terminal.Aliases, Role{ID: "terminal-extra", Alias: &Alias{Target: "does-not-exist"}})

	if _, err := ResolvedColors(palette); err == nil {
		t.Fatal("alias to a missing role was not rejected")
	}
}
