package fixture

import (
	"strings"
	"testing"
)

// requiredTerminalAgentStates is the contract this family must exercise. Every state
// must appear at least once across the scenes, and removing any single one must make
// the enumeration test fail by naming it.
func requiredTerminalAgentStates() []SemanticState {
	return []SemanticState{
		StateDefault, StateActive, StateInactive, StateFocused, StateSelected,
		StateMuted, StateSuccess, StateWarning, StateError, StateInfo,
	}
}

func TestTerminalAgentFixtureSetValidates(t *testing.T) {
	if err := TerminalAgentFixtureSet().Validate(); err != nil {
		t.Fatalf("Validate(): %v", err)
	}
}

func TestTerminalAgentFamilyExercisesRequiredStates(t *testing.T) {
	missing := missingRequiredStates(TerminalAgentFixtureSet())
	if len(missing) != 0 {
		t.Fatalf("missing required states: %v", missing)
	}
}

func TestTerminalAgentRequiredStateEnumerationReportsMissingState(t *testing.T) {
	for _, state := range requiredTerminalAgentStates() {
		t.Run(string(state), func(t *testing.T) {
			definition := TerminalAgentFixtureSet()
			replaceState(&definition, state, replacementFor(state))

			missing := missingRequiredStates(definition)
			if !containsState(missing, state) {
				t.Fatalf("enumeration did not report %q missing; missing=%v", state, missing)
			}
		})
	}
}

func TestTerminalAgentFamilyReferencesAllANSIRoles(t *testing.T) {
	roles := collectRunRoles(TerminalAgentFixtureSet())
	for _, role := range ansiRoleIDs {
		if _, ok := roles[role]; !ok {
			t.Fatalf("ANSI role %q is not referenced by any content run", role)
		}
	}
}

func TestTerminalAgentFamilyReferencesAllTerminalAliasRoles(t *testing.T) {
	roles := collectRunRoles(TerminalAgentFixtureSet())
	for _, role := range terminalAliasRoleIDs {
		if _, ok := roles[role]; !ok {
			t.Fatalf("terminal alias role %q is not referenced by any content run", role)
		}
	}
}

func TestTerminalAgentFamilyPlacesANSIBackgroundAgainstANSI0(t *testing.T) {
	for _, scene := range TerminalAgentFixtureSet().Scenes {
		for _, region := range scene.Regions {
			for _, block := range region.Blocks {
				for _, run := range block.Runs {
					if run.Foreground == RoleANSI0 && isANSIRole(run.Background) && run.Background != RoleANSI0 {
						return
					}
				}
			}
		}
	}
	t.Fatal("no content run places ANSI 0 text on a non-zero ANSI background")
}

func TestTerminalAgentFamilyValidatesAgainstMinimalPalette(t *testing.T) {
	if err := ValidatePair(TerminalAgentFixtureSet(), minimalPalette("palette-one")); err != nil {
		t.Fatalf("ValidatePair: %v", err)
	}
}

func TestTerminalAgentFamilySerializationIsDeterministic(t *testing.T) {
	encoded, err := MarshalIndent(TerminalAgentFixtureSet())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := readGolden(t, "terminal-agent.golden.json")
	if string(encoded) != want {
		t.Fatalf("serialized family differs from golden:\n%s", encoded)
	}
}

func TestTerminalAgentFamilyStructureIsPaletteIndependent(t *testing.T) {
	encoded, err := MarshalIndent(stripRoleReferences(TerminalAgentFixtureSet()))
	if err != nil {
		t.Fatalf("marshal structure: %v", err)
	}
	want := readGolden(t, "terminal-agent.structure.golden.json")
	if string(encoded) != want {
		t.Fatalf("family structure differs from golden:\n%s", encoded)
	}
}

func TestTerminalAgentFamilyUsesOnlyPlaceholderContent(t *testing.T) {
	for _, text := range collectText(TerminalAgentFixtureSet()) {
		for _, forbidden := range []string{"/home/", "user@"} {
			if strings.Contains(text, forbidden) {
				t.Fatalf("text %q contains forbidden placeholder %q", text, forbidden)
			}
		}
	}
}

func collectStates(definition FixtureSet) map[SemanticState]struct{} {
	states := make(map[SemanticState]struct{})
	for _, scene := range definition.Scenes {
		for _, region := range scene.Regions {
			states[region.State] = struct{}{}
			for _, block := range region.Blocks {
				states[block.State] = struct{}{}
				for _, run := range block.Runs {
					states[run.State] = struct{}{}
				}
			}
		}
	}
	return states
}

func missingRequiredStates(definition FixtureSet) []SemanticState {
	present := collectStates(definition)
	var missing []SemanticState
	for _, state := range requiredTerminalAgentStates() {
		if _, ok := present[state]; !ok {
			missing = append(missing, state)
		}
	}
	return missing
}

func replaceState(definition *FixtureSet, from, to SemanticState) {
	for sceneIndex := range definition.Scenes {
		scene := &definition.Scenes[sceneIndex]
		for regionIndex := range scene.Regions {
			region := &scene.Regions[regionIndex]
			if region.State == from {
				region.State = to
			}
			for blockIndex := range region.Blocks {
				block := &region.Blocks[blockIndex]
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

func replacementFor(state SemanticState) SemanticState {
	if state == StateDefault {
		return StateActive
	}
	return StateDefault
}

func collectRunRoles(definition FixtureSet) map[RoleID]struct{} {
	roles := make(map[RoleID]struct{})
	for _, scene := range definition.Scenes {
		for _, region := range scene.Regions {
			for _, block := range region.Blocks {
				for _, run := range block.Runs {
					if run.Foreground != "" {
						roles[run.Foreground] = struct{}{}
					}
					if run.Background != "" {
						roles[run.Background] = struct{}{}
					}
				}
			}
		}
	}
	return roles
}

func isANSIRole(role RoleID) bool {
	for _, candidate := range ansiRoleIDs {
		if candidate == role {
			return true
		}
	}
	return false
}

func stripRoleReferences(definition FixtureSet) FixtureSet {
	stripped := definition
	stripped.Scenes = make([]Scene, len(definition.Scenes))
	for sceneIndex, scene := range definition.Scenes {
		stripped.Scenes[sceneIndex] = scene
		regions := make([]Region, len(scene.Regions))
		for regionIndex, region := range scene.Regions {
			region.Background = ""
			region.Foreground = ""
			region.Border = ""
			blocks := make([]ContentBlock, len(region.Blocks))
			for blockIndex, block := range region.Blocks {
				block.Background = ""
				block.Foreground = ""
				block.Border = ""
				runs := make([]ContentRun, len(block.Runs))
				for runIndex, run := range block.Runs {
					run.Background = ""
					run.Foreground = ""
					runs[runIndex] = run
				}
				block.Runs = runs
				blocks[blockIndex] = block
			}
			region.Blocks = blocks
			regions[regionIndex] = region
		}
		stripped.Scenes[sceneIndex].Regions = regions
	}
	return stripped
}

func collectText(definition FixtureSet) []string {
	var texts []string
	for _, scene := range definition.Scenes {
		for _, region := range scene.Regions {
			for _, block := range region.Blocks {
				if block.Text != "" {
					texts = append(texts, block.Text)
				}
				for _, run := range block.Runs {
					texts = append(texts, run.Text)
				}
			}
		}
	}
	return texts
}

func containsState(states []SemanticState, wanted SemanticState) bool {
	for _, state := range states {
		if state == wanted {
			return true
		}
	}
	return false
}
