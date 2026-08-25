package fixture

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFreezeIsDeterministic(t *testing.T) {
	set := TerminalAgentFixtureSet()

	first, err := Freeze(set, FixtureSetVersionOne)
	if err != nil {
		t.Fatalf("first freeze: %v", err)
	}
	second, err := Freeze(set, FixtureSetVersionOne)
	if err != nil {
		t.Fatalf("second freeze: %v", err)
	}

	firstJSON, err := json.Marshal(first)
	if err != nil {
		t.Fatalf("marshal first: %v", err)
	}
	secondJSON, err := json.Marshal(second)
	if err != nil {
		t.Fatalf("marshal second: %v", err)
	}
	if string(firstJSON) != string(secondJSON) {
		t.Fatal("freeze is not deterministic: two freezes of the same content differ")
	}
}

func TestRegistryFreezeIsIdempotent(t *testing.T) {
	registry := NewRegistry()
	set := TerminalAgentFixtureSet()

	first, err := registry.Freeze(set, FixtureSetVersionOne)
	if err != nil {
		t.Fatalf("first freeze: %v", err)
	}
	second, err := registry.Freeze(set, FixtureSetVersionOne)
	if err != nil {
		t.Fatalf("second freeze: %v", err)
	}
	if first.Hash != second.Hash {
		t.Fatal("re-freezing the same content changed the hash")
	}
}

func TestRegistryFreezeRejectsVersionChange(t *testing.T) {
	registry := NewRegistry()
	set := TerminalAgentFixtureSet()
	if _, err := registry.Freeze(set, FixtureSetVersionOne); err != nil {
		t.Fatalf("initial freeze: %v", err)
	}

	mutated := TerminalAgentFixtureSet()
	mutated.Scenes[0].Regions[0].Blocks[1].Text = "changed"
	_, err := registry.Freeze(mutated, FixtureSetVersionOne)
	if err == nil {
		t.Fatal("freezing changed content under an existing version was accepted")
	}
	if !strings.Contains(err.Error(), "already frozen") {
		t.Fatalf("error %q does not explain the version conflict", err)
	}
}

func TestRegistryEntryCoverageRecordsEveryRequiredRole(t *testing.T) {
	entry, err := Freeze(TerminalAgentFixtureSet(), FixtureSetVersionOne)
	if err != nil {
		t.Fatalf("freeze: %v", err)
	}
	if len(entry.Coverage) != len(requiredCoverageRoles()) {
		t.Fatalf("coverage has %d entries, want %d", len(entry.Coverage), len(requiredCoverageRoles()))
	}
	for index, roleEntry := range entry.Coverage {
		if roleEntry.Role != requiredCoverageRoles()[index] {
			t.Fatalf("coverage entry %d role = %q, want %q", index, roleEntry.Role, requiredCoverageRoles()[index])
		}
	}
}

func TestRegistryEntryCoverageSpansAllLevels(t *testing.T) {
	entry, err := Freeze(TerminalAgentFixtureSet(), FixtureSetVersionOne)
	if err != nil {
		t.Fatalf("freeze: %v", err)
	}
	levels := make(map[ReferenceLevel]bool)
	for _, roleEntry := range entry.Coverage {
		for _, reference := range roleEntry.References {
			levels[reference.Level] = true
		}
	}
	for _, level := range []ReferenceLevel{LevelRegion, LevelBlock, LevelRun} {
		if !levels[level] {
			t.Fatalf("coverage does not span %s references", level)
		}
	}
}

func TestFrozenFixtureSetsSelectable(t *testing.T) {
	registry := FrozenFixtureSets()
	for _, set := range allFixtureSets() {
		entry, ok := registry.Entry(set.ID, FixtureSetVersionOne)
		if !ok {
			t.Fatalf("frozen entry for %q not found", set.ID)
		}
		if entry.Hash == "" {
			t.Fatalf("entry for %q has no hash", set.ID)
		}
		if len(entry.Scenes) != len(set.Scenes) {
			t.Fatalf("entry for %q has %d scenes, want %d", set.ID, len(entry.Scenes), len(set.Scenes))
		}
		if len(entry.Coverage) != len(requiredCoverageRoles()) {
			t.Fatalf("entry for %q has %d coverage entries, want %d", set.ID, len(entry.Coverage), len(requiredCoverageRoles()))
		}
	}
}
