package fixture

import "testing"

func allFixtureSets() []FixtureSet {
	return []FixtureSet{
		TerminalAgentFixtureSet(),
		CodeDiffFixtureSet(),
		DesktopShellFixtureSet(),
		ReadingMonitoringFixture(),
	}
}

func TestCanonicalHashIsDeterministic(t *testing.T) {
	for _, set := range allFixtureSets() {
		first, err := CanonicalHash(set)
		if err != nil {
			t.Fatalf("%s: %v", set.ID, err)
		}
		second, err := CanonicalHash(set)
		if err != nil {
			t.Fatalf("%s: %v", set.ID, err)
		}
		if first != second {
			t.Fatalf("%s: hash changed between identical calls", set.ID)
		}
	}
}

// Palette values cannot reach the canonical form, because Canonicalize takes only the
// fixture set and the palette is a separate document. Asserting that directly would be a
// tautology: the two hashes would come from the same unchanged input. What can actually
// break the property is a fixture that inlines a resolved color where a role reference
// belongs, which makes the structure carry a palette-specific value. That is the case
// TestCanonicalizeRejectsResolvedRoleValue covers, and it is where this obligation lives.

func TestCanonicalSceneHashesCoverEveryScene(t *testing.T) {
	for _, set := range allFixtureSets() {
		hashes, err := CanonicalSceneHashes(set)
		if err != nil {
			t.Fatalf("%s: %v", set.ID, err)
		}
		if len(hashes) != len(set.Scenes) {
			t.Fatalf("%s: got %d scene hashes, want %d", set.ID, len(hashes), len(set.Scenes))
		}
		for _, scene := range set.Scenes {
			if _, ok := hashes[scene.ID]; !ok {
				t.Fatalf("%s: scene %q has no canonical hash", set.ID, scene.ID)
			}
		}
	}
}

func TestEquivalentReportsIdenticalAndDivergent(t *testing.T) {
	set := TerminalAgentFixtureSet()
	same, err := Equivalent(set, set)
	if err != nil {
		t.Fatalf("Equivalent: %v", err)
	}
	if !same {
		t.Fatal("identical sets reported as divergent")
	}

	mutated := TerminalAgentFixtureSet()
	mutated.Scenes[0].Regions[0].Blocks[1].Text = "changed"
	different, err := Equivalent(set, mutated)
	if err != nil {
		t.Fatalf("Equivalent: %v", err)
	}
	if different {
		t.Fatal("divergent sets reported as equivalent")
	}
}

func TestCanonicalHashDetectsNonColorMutations(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(*FixtureSet)
		wantPath string
	}{
		{
			name: "content text",
			mutate: func(s *FixtureSet) {
				s.Scenes[0].Regions[0].Blocks[1].Text = "changed"
			},
			wantPath: "scenes[0].regions[0].blocks[1].text",
		},
		{
			name: "hierarchy block removed",
			mutate: func(s *FixtureSet) {
				blocks := s.Scenes[0].Regions[0].Blocks
				s.Scenes[0].Regions[0].Blocks = blocks[:len(blocks)-1]
			},
			wantPath: "scenes[0].regions[0].blocks[9]",
		},
		{
			name: "spacing block reordered",
			mutate: func(s *FixtureSet) {
				blocks := s.Scenes[0].Regions[0].Blocks
				blocks[0], blocks[1] = blocks[1], blocks[0]
			},
			wantPath: "scenes[0].regions[0].blocks[0].id",
		},
		{
			name: "state",
			mutate: func(s *FixtureSet) {
				s.Scenes[0].Regions[0].State = StateInactive
			},
			wantPath: "scenes[0].regions[0].state",
		},
		{
			name: "border role reference",
			mutate: func(s *FixtureSet) {
				s.Scenes[0].Regions[0].Border = RoleAccent
			},
			wantPath: "scenes[0].regions[0].border",
		},
		{
			name: "backdrop identity",
			mutate: func(s *FixtureSet) {
				s.Scenes[0].Regions[0].Background = RoleSurface
			},
			wantPath: "scenes[0].regions[0].background",
		},
		{
			name: "role-reference identity",
			mutate: func(s *FixtureSet) {
				s.Scenes[0].Regions[0].Foreground = RoleAccent
			},
			wantPath: "scenes[0].regions[0].foreground",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := TerminalAgentFixtureSet()
			mutated := TerminalAgentFixtureSet()
			test.mutate(&mutated)

			before, err := CanonicalHash(original)
			if err != nil {
				t.Fatalf("original hash: %v", err)
			}
			after, err := CanonicalHash(mutated)
			if err != nil {
				t.Fatalf("mutated hash: %v", err)
			}
			if before == after {
				t.Fatalf("mutation %q did not change the canonical hash", test.name)
			}

			paths, err := Divergence(original, mutated)
			if err != nil {
				t.Fatalf("Divergence: %v", err)
			}
			if !containsString(paths, test.wantPath) {
				t.Fatalf("divergence %v does not identify %q", paths, test.wantPath)
			}
		})
	}
}

func TestCanonicalizeRejectsResolvedRoleValue(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*FixtureSet)
	}{
		{
			name: "hex color in region background",
			mutate: func(s *FixtureSet) {
				s.Scenes[0].Regions[0].Background = "#ff0000"
			},
		},
		{
			name: "oklch literal in run foreground",
			mutate: func(s *FixtureSet) {
				s.Scenes[0].Regions[0].Blocks[0].Runs[0].Foreground = "oklch(0.5 0.1 200)"
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := TerminalAgentFixtureSet()
			test.mutate(&set)
			if _, err := Canonicalize(set); err == nil {
				t.Fatal("resolved color value in a role field was not rejected")
			}
		})
	}
}

func TestCanonicalGoldenForOneScenePerFamily(t *testing.T) {
	tests := []struct {
		name    string
		set     FixtureSet
		sceneID SceneID
		golden  string
	}{
		{"terminal-agent", TerminalAgentFixtureSet(), "terminal-shell", "canonical-terminal-agent.golden.json"},
		{"code-diff", CodeDiffFixtureSet(), "code-editor", "canonical-code-diff.golden.json"},
		{"desktop-shell", DesktopShellFixtureSet(), "desktop-shell-workspace", "canonical-desktop-shell.golden.json"},
		{"reading-monitoring", ReadingMonitoringFixture(), "reading-note", "canonical-reading-monitoring.golden.json"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := CanonicalScene(test.set, test.sceneID)
			if err != nil {
				t.Fatalf("canonical scene: %v", err)
			}
			want := readGolden(t, test.golden)
			if string(encoded) != want {
				t.Fatalf("canonical scene differs from golden:\n%s", encoded)
			}
		})
	}
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
