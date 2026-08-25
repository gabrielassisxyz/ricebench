package fixture

import (
	"strings"
	"testing"
)

func TestReadingMonitoringFixtureValidates(t *testing.T) {
	definition := ReadingMonitoringFixture()
	if err := definition.Validate(); err != nil {
		t.Fatalf("Validate(): %v", err)
	}
}

func TestReadingMonitoringRequiredStates(t *testing.T) {
	definition := ReadingMonitoringFixture()
	prose, proseOK := regionByID(definition, "prose")
	chrome, chromeOK := regionByID(definition, "note-chrome")
	monitor, monitorOK := regionByID(definition, "process-monitor")
	if !proseOK || !chromeOK || !monitorOK {
		t.Fatalf("regions missing: prose=%v chrome=%v monitor=%v", proseOK, chromeOK, monitorOK)
	}

	tests := []struct {
		name    string
		present func() bool
	}{
		{"heading", func() bool { return regionHasTextState(prose, StateActive) }},
		{"link", func() bool { return regionHasRun(prose, StateInfo, RoleAccent) }},
		{"inline-code", func() bool { return regionHasRunForeground(prose, RoleTextSecondary) }},
		{"fenced-code", func() bool { return regionHasKind(prose, ContentCode) }},
		{"muted-metadata", func() bool { return regionHasTextState(prose, StateMuted) }},
		{"active-chrome", func() bool { return regionHasTabState(chrome, StateActive) }},
		{"inactive-chrome", func() bool { return regionHasTabState(chrome, StateInactive) }},
		{"dense-table", func() bool { return hasRegionKind(definition, RegionTable) }},
		{"monitor-normal", func() bool { return regionHasStatusState(monitor, StateDefault) }},
		{"monitor-warning", func() bool { return regionHasStatusState(monitor, StateWarning) }},
		{"monitor-error", func() bool { return regionHasStatusState(monitor, StateError) }},
		{"monitor-success", func() bool { return regionHasStatusState(monitor, StateSuccess) }},
		{"monitor-muted", func() bool { return regionHasStatusState(monitor, StateMuted) }},
		{"monitor-info", func() bool { return regionHasStatusState(monitor, StateInfo) }},
		{"selection", func() bool { return hasSelectedRow(definition) }},
		{"scrolling", func() bool { return hasRegionKind(definition, RegionOverlay) }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.present() {
				t.Fatalf("required state %q is missing from the reading-monitoring family", test.name)
			}
		})
	}
}

func TestReadingMonitoringProseFloors(t *testing.T) {
	definition := ReadingMonitoringFixture()
	prose, ok := regionByID(definition, "prose")
	if !ok {
		t.Fatal("prose region missing")
	}

	var words int
	for _, block := range prose.Blocks {
		if block.Kind != ContentText || block.State != StateDefault {
			continue
		}
		line := blockText(block)
		if len(line) < 60 || len(line) > 90 {
			t.Errorf("prose line %q has length %d, want 60..90", block.ID, len(line))
		}
		words += len(strings.Fields(line))
	}
	if words < 400 {
		t.Errorf("prose carries %d words, want at least 400", words)
	}
}

func TestReadingMonitoringTableFloors(t *testing.T) {
	definition := ReadingMonitoringFixture()

	var rows int
	var selected, scrolled bool
	for _, scene := range definition.Scenes {
		for _, region := range scene.Regions {
			switch region.Kind {
			case RegionRow:
				if region.State != StateActive {
					rows++
				}
				if region.State == StateSelected {
					selected = true
				}
			case RegionOverlay:
				scrolled = true
			}
		}
	}

	if rows < 20 {
		t.Errorf("table carries %d data rows, want at least 20", rows)
	}
	if !selected {
		t.Error("table has no selected row")
	}
	if !scrolled {
		t.Error("table has no scrolled boundary")
	}
}

func TestReadingMonitoringMonitorAdjacency(t *testing.T) {
	definition := ReadingMonitoringFixture()
	monitor, ok := regionByID(definition, "process-monitor")
	if !ok {
		t.Fatal("monitor region missing")
	}

	want := []SemanticState{StateDefault, StateWarning, StateError, StateSuccess, StateMuted, StateInfo}
	var got []SemanticState
	for _, block := range monitor.Blocks {
		if block.Kind == ContentStatusItem {
			got = append(got, block.State)
		}
	}

	for _, state := range want {
		if !containsSemanticState(got, state) {
			t.Errorf("monitor region missing %q state", state)
		}
	}
	if !containsAdjacent(got, want) {
		t.Errorf("monitor states %v are not adjacent in %v", want, got)
	}
}

func TestReadingMonitoringConsumerFormatGuard(t *testing.T) {
	definition := ReadingMonitoringFixture()
	serialized, err := MarshalIndent(definition)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	lower := strings.ToLower(string(serialized))
	for _, name := range []string{"brave", "obsidian", "btop"} {
		if strings.Contains(lower, name) {
			t.Errorf("serialized family contains consumer product name %q", name)
		}
	}
}

func TestReadingMonitoringValidatePair(t *testing.T) {
	definition := ReadingMonitoringFixture()
	if err := ValidatePair(definition, minimalPalette("palette-one")); err != nil {
		t.Fatalf("ValidatePair: %v", err)
	}
}

func TestReadingMonitoringGolden(t *testing.T) {
	definition := ReadingMonitoringFixture()
	serialized, err := MarshalIndent(definition)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	want := readGolden(t, "reading-monitoring.golden.json")
	if string(serialized) != want {
		t.Fatalf("serialized family differs from golden:\n%s", serialized)
	}
}

func regionByID(set FixtureSet, id RegionID) (Region, bool) {
	for _, scene := range set.Scenes {
		for _, region := range scene.Regions {
			if region.ID == id {
				return region, true
			}
		}
	}
	return Region{}, false
}

func regionHasTextState(region Region, state SemanticState) bool {
	for _, block := range region.Blocks {
		if block.Kind == ContentText && block.State == state {
			return true
		}
	}
	return false
}

func regionHasRun(region Region, state SemanticState, foreground RoleID) bool {
	for _, block := range region.Blocks {
		for _, run := range block.Runs {
			if run.State == state && run.Foreground == foreground {
				return true
			}
		}
	}
	return false
}

func regionHasRunForeground(region Region, foreground RoleID) bool {
	for _, block := range region.Blocks {
		for _, run := range block.Runs {
			if run.Foreground == foreground {
				return true
			}
		}
	}
	return false
}

func regionHasKind(region Region, kind ContentKind) bool {
	for _, block := range region.Blocks {
		if block.Kind == kind {
			return true
		}
	}
	return false
}

func regionHasTabState(region Region, state SemanticState) bool {
	for _, block := range region.Blocks {
		if block.Kind == ContentTab && block.State == state {
			return true
		}
	}
	return false
}

func regionHasStatusState(region Region, state SemanticState) bool {
	for _, block := range region.Blocks {
		if block.Kind == ContentStatusItem && block.State == state {
			return true
		}
	}
	return false
}

func hasRegionKind(set FixtureSet, kind RegionKind) bool {
	for _, scene := range set.Scenes {
		for _, region := range scene.Regions {
			if region.Kind == kind {
				return true
			}
		}
	}
	return false
}

func hasSelectedRow(set FixtureSet) bool {
	for _, scene := range set.Scenes {
		for _, region := range scene.Regions {
			if region.Kind == RegionRow && region.State == StateSelected {
				return true
			}
		}
	}
	return false
}

func blockText(block ContentBlock) string {
	if block.Text != "" {
		return block.Text
	}
	var builder strings.Builder
	for _, run := range block.Runs {
		builder.WriteString(run.Text)
	}
	return builder.String()
}

func containsSemanticState(states []SemanticState, wanted SemanticState) bool {
	for _, state := range states {
		if state == wanted {
			return true
		}
	}
	return false
}

func containsAdjacent(states []SemanticState, wanted []SemanticState) bool {
	for start := 0; start+len(wanted) <= len(states); start++ {
		match := true
		for offset := range wanted {
			if states[start+offset] != wanted[offset] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
