package fixture

import (
	"slices"
	"strings"
	"testing"
)

func TestDesktopShellFixtureValidates(t *testing.T) {
	if err := DesktopShellFixtureSet().Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestDesktopShellValidatePair(t *testing.T) {
	if err := ValidatePair(DesktopShellFixtureSet(), minimalPalette("palette-one")); err != nil {
		t.Fatalf("ValidatePair: %v", err)
	}
}

func TestDesktopShellRequiredStates(t *testing.T) {
	missing := desktopShellMissingStates(DesktopShellFixtureSet())
	if len(missing) != 0 {
		t.Fatalf("missing required states: %s", strings.Join(missing, ", "))
	}
}

func TestDesktopShellRequiredStateRemoval(t *testing.T) {
	for _, state := range desktopShellRequiredStates() {
		t.Run(state.name, func(t *testing.T) {
			fixture := DesktopShellFixtureSet()
			state.remove(&fixture)
			missing := desktopShellMissingStates(fixture)
			if !slices.Contains(missing, state.name) {
				t.Fatalf("removing %q did not report it missing; missing=%v", state.name, missing)
			}
		})
	}
}

func TestDesktopShellConsumerFormatGuard(t *testing.T) {
	serialized, err := MarshalIndent(DesktopShellFixtureSet())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if violations := consumerFormatViolations(string(serialized)); len(violations) != 0 {
		t.Fatalf("clean family contains consumer-format names: %v", violations)
	}

	tests := []struct {
		name   string
		insert string
	}{
		{name: "waybar module key", insert: "modules-left"},
		{name: "hyprland option name", insert: "windowrule"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			violations := consumerFormatViolations(string(serialized) + test.insert)
			if !slices.Contains(violations, test.insert) {
				t.Fatalf("inserting %q did not report it; violations=%v", test.insert, violations)
			}
		})
	}
}

func TestDesktopShellWallpaperIsFlatSingleRole(t *testing.T) {
	wallpaper := findRegionByID(DesktopShellFixtureSet(), "wallpaper")
	if wallpaper == nil {
		t.Fatal("wallpaper region not found")
	}
	if wallpaper.Background != RoleBackground {
		t.Fatalf("wallpaper background = %q, want %q", wallpaper.Background, RoleBackground)
	}
	if wallpaper.Foreground != "" || wallpaper.Border != "" {
		t.Fatalf("wallpaper must be driven by a single role, got foreground=%q border=%q", wallpaper.Foreground, wallpaper.Border)
	}
	if len(wallpaper.Blocks) != 0 {
		t.Fatalf("wallpaper must be flat with no blocks, got %d", len(wallpaper.Blocks))
	}
}

func TestDesktopShellAdjacentSurfaces(t *testing.T) {
	if !hasAdjacentSurfaces(DesktopShellFixtureSet()) {
		t.Fatal("no two adjacent surfaces with distinct background roles")
	}
}

func TestDesktopShellFocusRingIsStructural(t *testing.T) {
	fixture := DesktopShellFixtureSet()
	focused := findRegionByID(fixture, "window-focused")
	if focused == nil {
		t.Fatal("focused window not found")
	}
	if !regionHasBlock(focused, ContentFocusRing) {
		t.Fatal("focused window has no focus-ring content block")
	}
	unfocused := findRegionByID(fixture, "window-unfocused")
	if unfocused == nil {
		t.Fatal("unfocused window not found")
	}
	if regionHasBlock(unfocused, ContentFocusRing) {
		t.Fatal("unfocused window must not carry a focus-ring")
	}
}

func TestDesktopShellGolden(t *testing.T) {
	serialized, err := MarshalIndent(DesktopShellFixtureSet())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := readGolden(t, "desktop-shell.golden.json")
	if string(serialized) != want {
		t.Fatalf("serialized family differs from golden:\n%s", serialized)
	}
}

type desktopShellRequiredState struct {
	name    string
	present func(FixtureSet) bool
	remove  func(*FixtureSet)
}

func desktopShellRequiredStates() []desktopShellRequiredState {
	return []desktopShellRequiredState{
		{
			name:    "active workspace",
			present: func(f FixtureSet) bool { return hasBlock(f, ContentStatusItem, StateActive) },
			remove:  func(f *FixtureSet) { mutateBlockState(f, ContentStatusItem, StateActive, StateDefault) },
		},
		{
			name:    "inactive workspace",
			present: func(f FixtureSet) bool { return hasBlock(f, ContentStatusItem, StateInactive) },
			remove:  func(f *FixtureSet) { mutateBlockState(f, ContentStatusItem, StateInactive, StateDefault) },
		},
		{
			name:    "focused window",
			present: func(f FixtureSet) bool { return hasRegion(f, RegionFrame, StateFocused) },
			remove:  func(f *FixtureSet) { mutateRegionState(f, RegionFrame, StateFocused, StateInactive) },
		},
		{
			name:    "unfocused window",
			present: func(f FixtureSet) bool { return hasRegion(f, RegionFrame, StateInactive) },
			remove:  func(f *FixtureSet) { mutateRegionState(f, RegionFrame, StateInactive, StateDefault) },
		},
		{
			name:    "status bar primary module",
			present: func(f FixtureSet) bool { return hasBlock(f, ContentStatusItem, StateDefault) },
			remove:  func(f *FixtureSet) { mutateBlockState(f, ContentStatusItem, StateDefault, StateMuted) },
		},
		{
			name:    "status bar muted module",
			present: func(f FixtureSet) bool { return hasBlock(f, ContentStatusItem, StateMuted) },
			remove:  func(f *FixtureSet) { mutateBlockState(f, ContentStatusItem, StateMuted, StateDefault) },
		},
		{
			name:    "launcher query",
			present: func(f FixtureSet) bool { return hasBlockID(f, "launcher-query") },
			remove:  func(f *FixtureSet) { removeBlockByID(f, "launcher-query") },
		},
		{
			name:    "launcher selected result",
			present: func(f FixtureSet) bool { return hasBlock(f, ContentListItem, StateSelected) },
			remove:  func(f *FixtureSet) { mutateBlockState(f, ContentListItem, StateSelected, StateDefault) },
		},
		{
			name:    "launcher secondary metadata",
			present: func(f FixtureSet) bool { return hasBlock(f, ContentText, StateMuted) },
			remove:  func(f *FixtureSet) { mutateBlockState(f, ContentText, StateMuted, StateDefault) },
		},
		{
			name:    "notification",
			present: func(f FixtureSet) bool { return hasRegionID(f, "notification") },
			remove:  func(f *FixtureSet) { removeRegionByID(f, "notification") },
		},
		{
			name:    "urgent state",
			present: func(f FixtureSet) bool { return hasAnyState(f, StateUrgent) },
			remove:  func(f *FixtureSet) { mutateAllState(f, StateUrgent, StateDefault) },
		},
		{
			name:    "transient overlay",
			present: func(f FixtureSet) bool { return hasRegionID(f, "transient-overlay") },
			remove:  func(f *FixtureSet) { removeRegionByID(f, "transient-overlay") },
		},
		{
			name:    "adjacent surfaces",
			present: hasAdjacentSurfaces,
			remove: func(f *FixtureSet) {
				for si := range f.Scenes {
					for ri := range f.Scenes[si].Regions {
						if f.Scenes[si].Regions[ri].ID == "window-focused" {
							f.Scenes[si].Regions[ri].Background = RoleSurface
						}
					}
				}
			},
		},
	}
}

func desktopShellMissingStates(f FixtureSet) []string {
	var missing []string
	for _, state := range desktopShellRequiredStates() {
		if !state.present(f) {
			missing = append(missing, state.name)
		}
	}
	return missing
}

var consumerFormatNames = []string{
	"hyprland", "waybar", "walker", "mako", "swayosd",
	"windowrule", "windowrulev2", "layerrule", "exec-once",
	"general:gaps_in", "decoration:rounding", "animations:enabled",
	"input:kb_layout", "monitor=", "bind=", "env=", "exec=",
	"xwayland", "blurls", "dwindle",
	"modules-left", "modules-right", "modules-center",
	"custom/", "on-click", "on-scroll-up", "on-scroll-down",
	"activation_mode",
	"max-visible", "default-timeout", "border-color", "background-color", "text-color", "group-by",
	"caps-lock",
}

func consumerFormatViolations(serialized string) []string {
	var violations []string
	for _, name := range consumerFormatNames {
		if strings.Contains(serialized, name) {
			violations = append(violations, name)
		}
	}
	return violations
}

func hasBlock(f FixtureSet, kind ContentKind, state SemanticState) bool {
	for _, scene := range f.Scenes {
		for _, region := range scene.Regions {
			for _, block := range region.Blocks {
				if block.Kind == kind && block.State == state {
					return true
				}
			}
		}
	}
	return false
}

func hasRegion(f FixtureSet, kind RegionKind, state SemanticState) bool {
	for _, scene := range f.Scenes {
		for _, region := range scene.Regions {
			if region.Kind == kind && region.State == state {
				return true
			}
		}
	}
	return false
}

func hasAnyState(f FixtureSet, state SemanticState) bool {
	for _, scene := range f.Scenes {
		for _, region := range scene.Regions {
			if region.State == state {
				return true
			}
			for _, block := range region.Blocks {
				if block.State == state {
					return true
				}
				for _, run := range block.Runs {
					if run.State == state {
						return true
					}
				}
			}
		}
	}
	return false
}

func hasBlockID(f FixtureSet, id ContentBlockID) bool {
	for _, scene := range f.Scenes {
		for _, region := range scene.Regions {
			for _, block := range region.Blocks {
				if block.ID == id {
					return true
				}
			}
		}
	}
	return false
}

func hasRegionID(f FixtureSet, id RegionID) bool {
	for _, scene := range f.Scenes {
		for _, region := range scene.Regions {
			if region.ID == id {
				return true
			}
		}
	}
	return false
}

func hasAdjacentSurfaces(f FixtureSet) bool {
	for _, scene := range f.Scenes {
		for i := 0; i+1 < len(scene.Regions); i++ {
			a, b := scene.Regions[i].Background, scene.Regions[i+1].Background
			if (a == RoleSurface && b == RoleSurfaceElevated) || (a == RoleSurfaceElevated && b == RoleSurface) {
				return true
			}
		}
	}
	return false
}

func mutateBlockState(f *FixtureSet, kind ContentKind, from, to SemanticState) {
	for si := range f.Scenes {
		for ri := range f.Scenes[si].Regions {
			for bi := range f.Scenes[si].Regions[ri].Blocks {
				block := &f.Scenes[si].Regions[ri].Blocks[bi]
				if block.Kind == kind && block.State == from {
					block.State = to
				}
			}
		}
	}
}

func mutateRegionState(f *FixtureSet, kind RegionKind, from, to SemanticState) {
	for si := range f.Scenes {
		for ri := range f.Scenes[si].Regions {
			region := &f.Scenes[si].Regions[ri]
			if region.Kind == kind && region.State == from {
				region.State = to
			}
		}
	}
}

func mutateAllState(f *FixtureSet, from, to SemanticState) {
	for si := range f.Scenes {
		for ri := range f.Scenes[si].Regions {
			region := &f.Scenes[si].Regions[ri]
			if region.State == from {
				region.State = to
			}
			for bi := range region.Blocks {
				block := &region.Blocks[bi]
				if block.State == from {
					block.State = to
				}
				for runIndex := range block.Runs {
					if block.Runs[runIndex].State == from {
						block.Runs[runIndex].State = to
					}
				}
			}
		}
	}
}

func removeBlockByID(f *FixtureSet, id ContentBlockID) {
	for si := range f.Scenes {
		for ri := range f.Scenes[si].Regions {
			kept := make([]ContentBlock, 0, len(f.Scenes[si].Regions[ri].Blocks))
			for _, block := range f.Scenes[si].Regions[ri].Blocks {
				if block.ID != id {
					kept = append(kept, block)
				}
			}
			f.Scenes[si].Regions[ri].Blocks = kept
		}
	}
}

func removeRegionByID(f *FixtureSet, id RegionID) {
	for si := range f.Scenes {
		kept := make([]Region, 0, len(f.Scenes[si].Regions))
		for _, region := range f.Scenes[si].Regions {
			if region.ID != id {
				kept = append(kept, region)
			}
		}
		f.Scenes[si].Regions = kept
	}
}

func findRegionByID(f FixtureSet, id RegionID) *Region {
	for si := range f.Scenes {
		for ri := range f.Scenes[si].Regions {
			if f.Scenes[si].Regions[ri].ID == id {
				return &f.Scenes[si].Regions[ri]
			}
		}
	}
	return nil
}

func regionHasBlock(region *Region, kind ContentKind) bool {
	for _, block := range region.Blocks {
		if block.Kind == kind {
			return true
		}
	}
	return false
}
