package fixture

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidatePairAcceptsOneFixtureWithMultiplePalettes(t *testing.T) {
	definition := minimalFixture()
	first := minimalPalette("palette-one")
	second := minimalPalette("palette-two")

	before, err := json.Marshal(definition)
	if err != nil {
		t.Fatalf("marshal fixture before validation: %v", err)
	}

	for _, palette := range []Palette{first, second} {
		if err := ValidatePair(definition, palette); err != nil {
			t.Fatalf("ValidatePair(%q): %v", palette.ID, err)
		}
	}

	after, err := json.Marshal(definition)
	if err != nil {
		t.Fatalf("marshal fixture after validation: %v", err)
	}
	if string(after) != string(before) {
		t.Fatal("pair validation mutated the fixture definition")
	}
}

func TestFixtureValidationRejectsInvalidContracts(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*FixtureSet)
		wantErr string
	}{
		{
			name: "unknown schema version",
			mutate: func(definition *FixtureSet) {
				definition.SchemaVersion = "2"
			},
			wantErr: `schemaVersion: got "2", want "1"`,
		},
		{
			name: "invalid fixture identity",
			mutate: func(definition *FixtureSet) {
				definition.ID = "Not Stable"
			},
			wantErr: `id: "Not Stable" is not a stable ID`,
		},
		{
			name: "duplicate scene identity",
			mutate: func(definition *FixtureSet) {
				definition.Scenes = append(definition.Scenes, definition.Scenes[0])
			},
			wantErr: `scenes[1].id: duplicate "terminal-shell"`,
		},
		{
			name: "unknown scene family",
			mutate: func(definition *FixtureSet) {
				definition.Scenes[0].Family = "marketplace"
			},
			wantErr: `scenes[0].family: unknown "marketplace"`,
		},
		{
			name: "duplicate region identity",
			mutate: func(definition *FixtureSet) {
				definition.Scenes[0].Regions = append(
					definition.Scenes[0].Regions,
					definition.Scenes[0].Regions[0],
				)
			},
			wantErr: `scenes[0].regions[1].id: duplicate "shell-output"`,
		},
		{
			name: "duplicate content block identity",
			mutate: func(definition *FixtureSet) {
				blocks := definition.Scenes[0].Regions[0].Blocks
				definition.Scenes[0].Regions[0].Blocks = append(blocks, blocks[0])
			},
			wantErr: `blocks[1].id: duplicate "command-result"`,
		},
		{
			name: "unknown semantic state",
			mutate: func(definition *FixtureSet) {
				definition.Scenes[0].Regions[0].Blocks[0].State = "sparkling"
			},
			wantErr: `state: unknown "sparkling"`,
		},
		{
			name: "ambiguous block content",
			mutate: func(definition *FixtureSet) {
				definition.Scenes[0].Regions[0].Blocks[0].Runs = []ContentRun{
					{Text: "also ok", State: StateSuccess, Foreground: RoleSuccess},
				}
			},
			wantErr: `text and runs are mutually exclusive`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			definition := minimalFixture()
			test.mutate(&definition)
			assertErrorContains(t, definition.Validate(), test.wantErr)
		})
	}
}

func TestPaletteValidationRejectsInvalidContracts(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Palette)
		wantErr string
	}{
		{
			name: "unknown schema version",
			mutate: func(palette *Palette) {
				palette.SchemaVersion = "2"
			},
			wantErr: `schemaVersion: got "2", want "1"`,
		},
		{
			name: "duplicate role",
			mutate: func(palette *Palette) {
				palette.SemanticCore[1].ID = palette.SemanticCore[0].ID
			},
			wantErr: `semanticCore[1].id: duplicate "background"`,
		},
		{
			name: "missing required role",
			mutate: func(palette *Palette) {
				palette.SemanticCore = palette.SemanticCore[1:]
			},
			wantErr: `semanticCore: missing required role "background"`,
		},
		{
			name: "authored role with alias",
			mutate: func(palette *Palette) {
				palette.SemanticCore[0].Alias = &Alias{Target: RoleForeground}
			},
			wantErr: `semanticCore[0]: role "background" must contain exactly one of value or alias`,
		},
		{
			name: "alias target does not exist",
			mutate: func(palette *Palette) {
				palette.Terminal.Aliases[0].Alias.Target = "missing-role"
			},
			wantErr: `terminal.aliases[0].alias.target: role "missing-role" does not exist`,
		},
		{
			name: "alias cycle",
			mutate: func(palette *Palette) {
				palette.Terminal.Aliases[0].Alias.Target = RoleTerminalForeground
				palette.Terminal.Aliases[1].Alias.Target = RoleTerminalBackground
			},
			wantErr: `alias cycle: terminal-background -> terminal-foreground -> terminal-background`,
		},
		{
			name: "terminal alias resolves outside semantic core",
			mutate: func(palette *Palette) {
				palette.Terminal.Aliases[0].Alias.Target = RoleANSI0
			},
			wantErr: `terminal alias "terminal-background" resolves to non-semantic role "ansi-0"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			palette := minimalPalette("palette-one")
			test.mutate(&palette)
			assertErrorContains(t, palette.Validate(), test.wantErr)
		})
	}
}

func TestValidatePairRejectsMissingSceneRole(t *testing.T) {
	definition := minimalFixture()
	definition.Scenes[0].Regions[0].Blocks[0].Foreground = "fixture-only-role"

	err := ValidatePair(definition, minimalPalette("palette-one"))
	assertErrorContains(t, err, `role "fixture-only-role" does not exist in palette "palette-one"`)
}

func TestValidatePairRejectsMissingRunRole(t *testing.T) {
	definition := minimalFixture()
	block := &definition.Scenes[0].Regions[0].Blocks[0]
	block.Text = ""
	block.Runs = []ContentRun{
		{Text: "ok", State: StateSuccess, Foreground: "fixture-only-role"},
	}

	err := ValidatePair(definition, minimalPalette("palette-one"))
	assertErrorContains(t, err, `runs[0]: role "fixture-only-role" does not exist in palette "palette-one"`)
}

func TestMarshalIndentIsDeterministicAndHumanDiffable(t *testing.T) {
	document := minimalFixture()

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
	if !strings.HasSuffix(string(first), "\n") {
		t.Fatal("serialized document must end with a newline")
	}

	want := readGolden(t, "minimal.golden.json")
	if string(first) != want {
		t.Fatalf("serialized document differs from golden:\n%s", first)
	}
}

func assertErrorContains(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Fatalf("got nil error, want one containing %q", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error %q does not contain %q", err, want)
	}
}
