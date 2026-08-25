package fixture

import (
	"slices"
	"strings"
	"testing"
)

// codeDiffRequiredStates is the contract the code-diff family must satisfy. The
// enumeration test below fails by name when any one of these is absent.
var codeDiffRequiredStates = []string{
	"keyword", "type", "function", "string", "number", "comment", "punctuation",
	"diagnostic-error", "diagnostic-warning", "inline-annotation",
	"added-line", "removed-line", "modified-line",
	"selected-tab", "inactive-tab",
	"search-match", "active-line", "focused-control",
	"dense-code",
}

func TestCodeDiffFamilyValidates(t *testing.T) {
	set := CodeDiffFixtureSet()
	if err := set.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if err := ValidatePair(set, minimalPalette("palette-one")); err != nil {
		t.Fatalf("ValidatePair: %v", err)
	}
}

func TestCodeDiffFamilyHasAllRequiredStates(t *testing.T) {
	missing := codeDiffMissingStates(CodeDiffFixtureSet())
	if len(missing) != 0 {
		t.Fatalf("code-diff family missing required states: %v", missing)
	}
}

func TestCodeDiffFamilyReportsMissingState(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*FixtureSet)
		wantErr string
	}{
		{name: "keyword", mutate: func(s *FixtureSet) { stripRunRole(s, codeKeywordRole) }, wantErr: "keyword"},
		{name: "type", mutate: func(s *FixtureSet) { stripRunRole(s, codeTypeRole) }, wantErr: "type"},
		{name: "function", mutate: func(s *FixtureSet) { stripRunRole(s, codeFunctionRole) }, wantErr: "function"},
		{name: "string", mutate: func(s *FixtureSet) { stripRunRole(s, codeStringRole) }, wantErr: "string"},
		{name: "number", mutate: func(s *FixtureSet) { stripRunRole(s, codeNumberRole) }, wantErr: "number"},
		{name: "comment", mutate: func(s *FixtureSet) { stripRunRole(s, codeCommentRole) }, wantErr: "comment"},
		{name: "punctuation", mutate: func(s *FixtureSet) { stripRunRole(s, codePunctuationRole) }, wantErr: "punctuation"},
		{name: "diagnostic-error", mutate: func(s *FixtureSet) { stripRunState(s, StateError) }, wantErr: "diagnostic-error"},
		{name: "diagnostic-warning", mutate: func(s *FixtureSet) { stripRunState(s, StateWarning) }, wantErr: "diagnostic-warning"},
		{name: "inline-annotation", mutate: stripAnnotationText, wantErr: "inline-annotation"},
		{name: "added-line", mutate: func(s *FixtureSet) { stripBlockState(s, StateAdded) }, wantErr: "added-line"},
		{name: "removed-line", mutate: func(s *FixtureSet) { stripBlockState(s, StateRemoved) }, wantErr: "removed-line"},
		{name: "modified-line", mutate: func(s *FixtureSet) { stripBlockState(s, StateModified) }, wantErr: "modified-line"},
		{name: "selected-tab", mutate: func(s *FixtureSet) { stripTabState(s, StateSelected) }, wantErr: "selected-tab"},
		{name: "inactive-tab", mutate: func(s *FixtureSet) { stripTabState(s, StateInactive) }, wantErr: "inactive-tab"},
		{name: "search-match", mutate: func(s *FixtureSet) { stripRunState(s, StateSearchMatch) }, wantErr: "search-match"},
		{name: "active-line", mutate: func(s *FixtureSet) { stripBlockState(s, StateActive) }, wantErr: "active-line"},
		{name: "focused-control", mutate: func(s *FixtureSet) { stripBlockState(s, StateFocused) }, wantErr: "focused-control"},
		{name: "dense-code", mutate: truncateCode, wantErr: "dense-code"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := CodeDiffFixtureSet()
			test.mutate(&set)
			missing := codeDiffMissingStates(set)
			if !slices.Contains(missing, test.wantErr) {
				t.Fatalf("missing states %v do not contain %q", missing, test.wantErr)
			}
		})
	}
}

func TestCodeDiffGutterMarkersDistinguishableUnderGrayscale(t *testing.T) {
	set := CodeDiffFixtureSet()
	// minimalPalette authors every role as the same gray, so any distinction the
	// test still observes cannot come from hue.
	if err := ValidatePair(set, minimalPalette("grayscale")); err != nil {
		t.Fatalf("grayscale pair: %v", err)
	}

	markers := codeDiffGutterMarkers(set)
	for _, state := range []SemanticState{StateAdded, StateRemoved, StateModified} {
		if markers[state] == "" {
			t.Fatalf("state %q has no gutter marker", state)
		}
	}
	if markers[StateAdded] == markers[StateRemoved] ||
		markers[StateAdded] == markers[StateModified] ||
		markers[StateRemoved] == markers[StateModified] {
		t.Fatalf("gutter markers must be distinct: %v", markers)
	}
}

func TestCodeDiffDiagnosticsDistinguishableWithoutHue(t *testing.T) {
	annotations := codeDiffAnnotations(CodeDiffFixtureSet())

	if !strings.Contains(annotations[StateError], "error") {
		t.Fatalf("error annotation %q does not name its severity", annotations[StateError])
	}
	if !strings.Contains(annotations[StateWarning], "warning") {
		t.Fatalf("warning annotation %q does not name its severity", annotations[StateWarning])
	}
	if annotations[StateError] == annotations[StateWarning] {
		t.Fatalf("error and warning annotations must differ: %q", annotations[StateError])
	}
}

func TestCodeDiffMarshalIndentIsDeterministic(t *testing.T) {
	document := CodeDiffFixtureSet()

	first, err := MarshalIndent(document)
	if err != nil {
		t.Fatalf("first marshal: %v", err)
	}
	second, err := MarshalIndent(document)
	if err != nil {
		t.Fatalf("second marshal: %v", err)
	}
	if string(first) != string(second) {
		t.Fatal("serialization changed between identical calls")
	}

	want := readGolden(t, "code-diff.golden.json")
	if string(first) != want {
		t.Fatalf("serialized document differs from golden:\n%s", first)
	}
}

func TestCodeDiffStructureSnapshot(t *testing.T) {
	document := stripRoleReferences(CodeDiffFixtureSet())
	encoded, err := MarshalIndent(document)
	if err != nil {
		t.Fatalf("marshal stripped structure: %v", err)
	}

	want := readGolden(t, "code-diff.structure.golden.json")
	if string(encoded) != want {
		t.Fatalf("stripped structure differs from snapshot:\n%s", encoded)
	}
}

// codeDiffPresentStates reports which required states the family currently carries.
func codeDiffPresentStates(set FixtureSet) map[string]bool {
	present := make(map[string]bool)
	codeLineCount := 0
	for _, scene := range set.Scenes {
		for _, region := range scene.Regions {
			for _, block := range region.Blocks {
				if block.Kind == ContentCode && len(block.Runs) > 0 {
					codeLineCount++
				}
				switch {
				case block.Kind == ContentTab && block.State == StateSelected:
					present["selected-tab"] = true
				case block.Kind == ContentTab && block.State == StateInactive:
					present["inactive-tab"] = true
				case block.Kind == ContentCode && block.State == StateAdded && len(block.Runs) > 0:
					present["added-line"] = true
				case block.Kind == ContentCode && block.State == StateRemoved && len(block.Runs) > 0:
					present["removed-line"] = true
				case block.Kind == ContentCode && block.State == StateModified && len(block.Runs) > 0:
					present["modified-line"] = true
				case block.Kind == ContentCode && block.State == StateActive && len(block.Runs) > 0:
					present["active-line"] = true
				case block.State == StateFocused:
					present["focused-control"] = true
				case block.Kind == ContentText && isAnnotation(block.Text):
					present["inline-annotation"] = true
				}
				for _, run := range block.Runs {
					switch run.Foreground {
					case codeKeywordRole:
						present["keyword"] = true
					case codeTypeRole:
						present["type"] = true
					case codeFunctionRole:
						present["function"] = true
					case codeStringRole:
						present["string"] = true
					case codeNumberRole:
						present["number"] = true
					case codeCommentRole:
						present["comment"] = true
					case codePunctuationRole:
						present["punctuation"] = true
					}
					switch run.State {
					case StateError:
						present["diagnostic-error"] = true
					case StateWarning:
						present["diagnostic-warning"] = true
					case StateSearchMatch:
						present["search-match"] = true
					}
				}
			}
		}
	}
	if codeLineCount >= 40 {
		present["dense-code"] = true
	}
	return present
}

func codeDiffMissingStates(set FixtureSet) []string {
	present := codeDiffPresentStates(set)
	var missing []string
	for _, state := range codeDiffRequiredStates {
		if !present[state] {
			missing = append(missing, state)
		}
	}
	return missing
}

// codeDiffGutterMarkers returns the gutter text for each diff state. The markers are
// text, so they carry meaning independently of hue.
func codeDiffGutterMarkers(set FixtureSet) map[SemanticState]string {
	markers := make(map[SemanticState]string)
	for _, scene := range set.Scenes {
		for _, region := range scene.Regions {
			for _, block := range region.Blocks {
				if block.Text == "" {
					continue
				}
				switch block.State {
				case StateAdded, StateRemoved, StateModified:
					markers[block.State] = block.Text
				}
			}
		}
	}
	return markers
}

// codeDiffAnnotations returns the inline annotation text keyed by severity.
func codeDiffAnnotations(set FixtureSet) map[SemanticState]string {
	annotations := make(map[SemanticState]string)
	for _, scene := range set.Scenes {
		for _, region := range scene.Regions {
			for _, block := range region.Blocks {
				if block.Kind != ContentText || !isAnnotation(block.Text) {
					continue
				}
				annotations[block.State] = block.Text
			}
		}
	}
	return annotations
}

func isAnnotation(text string) bool {
	return strings.Contains(text, "error:") || strings.Contains(text, "warning:")
}

// stripRoleReferences zeroes every role reference so the remaining structure is the
// local half of the color-only equivalence proof completed in P4.7.
func stripRoleReferences(set FixtureSet) FixtureSet {
	for si := range set.Scenes {
		for ri := range set.Scenes[si].Regions {
			region := &set.Scenes[si].Regions[ri]
			region.Background = ""
			region.Foreground = ""
			region.Border = ""
			for bi := range region.Blocks {
				block := &region.Blocks[bi]
				block.Background = ""
				block.Foreground = ""
				block.Border = ""
				for runi := range block.Runs {
					block.Runs[runi].Background = ""
					block.Runs[runi].Foreground = ""
				}
			}
		}
	}
	return set
}

func stripRunRole(set *FixtureSet, role RoleID) {
	for si := range set.Scenes {
		for ri := range set.Scenes[si].Regions {
			for bi := range set.Scenes[si].Regions[ri].Blocks {
				for runi := range set.Scenes[si].Regions[ri].Blocks[bi].Runs {
					if set.Scenes[si].Regions[ri].Blocks[bi].Runs[runi].Foreground == role {
						set.Scenes[si].Regions[ri].Blocks[bi].Runs[runi].Foreground = RoleForeground
					}
				}
			}
		}
	}
}

func stripRunState(set *FixtureSet, state SemanticState) {
	for si := range set.Scenes {
		for ri := range set.Scenes[si].Regions {
			for bi := range set.Scenes[si].Regions[ri].Blocks {
				for runi := range set.Scenes[si].Regions[ri].Blocks[bi].Runs {
					if set.Scenes[si].Regions[ri].Blocks[bi].Runs[runi].State == state {
						set.Scenes[si].Regions[ri].Blocks[bi].Runs[runi].State = StateDefault
					}
				}
			}
		}
	}
}

func stripBlockState(set *FixtureSet, state SemanticState) {
	for si := range set.Scenes {
		for ri := range set.Scenes[si].Regions {
			for bi := range set.Scenes[si].Regions[ri].Blocks {
				if set.Scenes[si].Regions[ri].Blocks[bi].State == state {
					set.Scenes[si].Regions[ri].Blocks[bi].State = StateDefault
				}
			}
		}
	}
}

func stripTabState(set *FixtureSet, state SemanticState) {
	for si := range set.Scenes {
		for ri := range set.Scenes[si].Regions {
			for bi := range set.Scenes[si].Regions[ri].Blocks {
				block := &set.Scenes[si].Regions[ri].Blocks[bi]
				if block.Kind == ContentTab && block.State == state {
					block.State = StateDefault
				}
			}
		}
	}
}

func stripAnnotationText(set *FixtureSet) {
	for si := range set.Scenes {
		for ri := range set.Scenes[si].Regions {
			for bi := range set.Scenes[si].Regions[ri].Blocks {
				block := &set.Scenes[si].Regions[ri].Blocks[bi]
				if block.Kind == ContentText && isAnnotation(block.Text) {
					block.Text = "note"
				}
			}
		}
	}
}

func truncateCode(set *FixtureSet) {
	for si := range set.Scenes {
		for ri := range set.Scenes[si].Regions {
			blocks := set.Scenes[si].Regions[ri].Blocks
			kept := blocks[:0]
			for _, block := range blocks {
				if block.Kind == ContentCode && len(block.Runs) > 0 {
					continue
				}
				kept = append(kept, block)
			}
			set.Scenes[si].Regions[ri].Blocks = kept
		}
	}
}
