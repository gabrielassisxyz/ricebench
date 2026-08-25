package fixture

import "strconv"

// Syntax token roles for the code-diff family. Each token class maps to a distinct
// ANSI slot so syntax density is exposed rather than flattened into one foreground.
const (
	codeKeywordRole     RoleID = RoleANSI4
	codeTypeRole        RoleID = RoleANSI6
	codeFunctionRole    RoleID = RoleANSI2
	codeStringRole      RoleID = RoleANSI3
	codeNumberRole      RoleID = RoleANSI5
	codeCommentRole     RoleID = RoleANSI8
	codePunctuationRole RoleID = RoleANSI7
)

// CodeDiffFixtureSet returns the code-diff scene family: a dense code editor with
// syntax highlighting, diagnostics, and a diff view. Every candidate renders from
// this one definition, so a comparison can only differ in color.
func CodeDiffFixtureSet() FixtureSet {
	return FixtureSet{
		SchemaVersion: FixtureSchemaVersion,
		ID:            "fixture-code-diff",
		Scenes: []Scene{
			codeEditorScene(),
			diffViewScene(),
		},
	}
}

func codeEditorScene() Scene {
	return Scene{
		ID:     "code-editor",
		Family: FamilyCodeDiff,
		Regions: []Region{
			{
				ID:         "editor-frame",
				Kind:       RegionFrame,
				State:      StateActive,
				Background: RoleBackground,
				Foreground: RoleForeground,
			},
			{
				ID:    "tab-bar",
				Kind:  RegionTabs,
				State: StateDefault,
				Blocks: []ContentBlock{
					{ID: "tab-selected", Kind: ContentTab, State: StateSelected, Text: "src/lexer.go", Foreground: RoleAccent},
					{ID: "tab-inactive", Kind: ContentTab, State: StateInactive, Text: "src/parser.go", Foreground: RoleTextMuted},
					{ID: "tab-inactive-2", Kind: ContentTab, State: StateInactive, Text: "src/ast.go", Foreground: RoleTextMuted},
				},
			},
			{
				ID:     "search-bar",
				Kind:   RegionRow,
				State:  StateFocused,
				Border: RoleFocus,
				Blocks: []ContentBlock{
					{ID: "search-label", Kind: ContentText, State: StateDefault, Text: "search:", Foreground: RoleTextMuted},
					{ID: "search-input", Kind: ContentText, State: StateFocused, Text: "Token", Foreground: RoleForeground},
					{ID: "search-count", Kind: ContentStatusItem, State: StateSearchMatch, Text: "3 matches", Foreground: RoleAccent},
				},
			},
			{
				ID:         "gutter",
				Kind:       RegionColumn,
				State:      StateDefault,
				Background: RoleSurface,
				Blocks:     lineNumberBlocks(codeEditorLineCount),
			},
			{
				ID:         "code",
				Kind:       RegionColumn,
				State:      StateDefault,
				Background: RoleSurface,
				Blocks:     codeEditorCodeLines(),
			},
			{
				ID:         "status-bar",
				Kind:       RegionStatus,
				State:      StateDefault,
				Background: RoleSurfaceElevated,
				Blocks: []ContentBlock{
					{ID: "status-diagnostics", Kind: ContentStatusItem, State: StateError, Text: "1 error, 1 warning", Foreground: RoleError},
					{ID: "status-mode", Kind: ContentStatusItem, State: StateDefault, Text: "NORMAL", Foreground: RoleTextSecondary},
				},
			},
		},
	}
}

func diffViewScene() Scene {
	return Scene{
		ID:     "diff-view",
		Family: FamilyCodeDiff,
		Regions: []Region{
			{
				ID:         "diff-frame",
				Kind:       RegionFrame,
				State:      StateActive,
				Background: RoleBackground,
				Foreground: RoleForeground,
			},
			{
				ID:    "diff-tabs",
				Kind:  RegionTabs,
				State: StateDefault,
				Blocks: []ContentBlock{
					{ID: "diff-tab-selected", Kind: ContentTab, State: StateSelected, Text: "src/lexer.go", Foreground: RoleAccent},
					{ID: "diff-tab-inactive", Kind: ContentTab, State: StateInactive, Text: "src/parser.go", Foreground: RoleTextMuted},
				},
			},
			{
				ID:    "diff-hunk",
				Kind:  RegionRow,
				State: StateDefault,
				Blocks: []ContentBlock{
					{ID: "hunk-header", Kind: ContentText, State: StateDefault, Text: "@@ -4,3 +4,3 @@ type Token struct", Foreground: RoleTextSecondary},
				},
			},
			{
				ID:         "diff-gutter",
				Kind:       RegionColumn,
				State:      StateDefault,
				Background: RoleSurface,
				Blocks: []ContentBlock{
					{ID: "diff-gutter-4", Kind: ContentCode, State: StateDefault, Text: "4", Foreground: RoleTextMuted},
					{ID: "diff-gutter-removed", Kind: ContentCode, State: StateRemoved, Text: "-", Foreground: RoleDiffRemove},
					{ID: "diff-gutter-added", Kind: ContentCode, State: StateAdded, Text: "+", Foreground: RoleDiffAdd},
					{ID: "diff-gutter-modified", Kind: ContentCode, State: StateModified, Text: "~", Foreground: RoleDiffModify},
					{ID: "diff-gutter-5", Kind: ContentCode, State: StateDefault, Text: "5", Foreground: RoleTextMuted},
				},
			},
			{
				ID:         "diff-code",
				Kind:       RegionColumn,
				State:      StateDefault,
				Background: RoleSurface,
				Blocks: []ContentBlock{
					codeLine("diff-code-4", StateDefault, ws("    "), ty("Kind"), ws("    "), ty("string")),
					codeLine("diff-code-removed", StateRemoved, ws("    "), ty("Literal"), ws(" "), ty("string")),
					codeLine("diff-code-added", StateAdded, ws("    "), ty("Literal"), ws(" "), punc("["), punc("]"), ty("byte")),
					codeLine("diff-code-modified", StateModified, ws("    "), ty("Line"), ws("    "), ty("uint")),
					codeLine("diff-code-5", StateDefault, ws("    "), ty("Col"), ws("     "), ty("int")),
				},
			},
			{
				ID:         "diff-status",
				Kind:       RegionStatus,
				State:      StateDefault,
				Background: RoleSurfaceElevated,
				Blocks: []ContentBlock{
					{ID: "diff-status-summary", Kind: ContentStatusItem, State: StateDefault, Text: "3 changes", Foreground: RoleTextSecondary},
				},
			},
		},
	}
}

// codeEditorLineCount is the number of code lines in the editor scene. It is kept
// above the forty-line floor so hierarchy and fatigue are judgeable rather than
// inferred from a snippet.
const codeEditorLineCount = 52

func codeEditorCodeLines() []ContentBlock {
	return []ContentBlock{
		codeLine("code-1", StateDefault, kw("package"), ws(" "), ty("lexer")),
		codeLine("code-2", StateDefault, kw("import"), ws(" "), str(`"fmt"`)),
		codeLine("code-3", StateDefault, com("// Token is one lexical unit produced by the scanner.")),
		codeLine("code-4", StateDefault, kw("type"), ws(" "), ty("Token").marked(StateSearchMatch).bg(RoleSelectionBackground), ws(" "), kw("struct"), ws(" "), punc("{")),
		codeLine("code-5", StateDefault, ws("    "), ty("Kind"), ws("    "), ty("string")),
		codeLine("code-6", StateDefault, ws("    "), ty("Literal"), ws(" "), ty("string")),
		codeLine("code-7", StateDefault, ws("    "), ty("Line"), ws("    "), ty("int")),
		codeLine("code-8", StateDefault, punc("}")),
		codeLine("code-9", StateDefault, com("// Scanner walks source text and emits tokens.")),
		codeLine("code-10", StateDefault, kw("type"), ws(" "), ty("Scanner"), ws(" "), kw("struct"), ws(" "), punc("{")),
		codeLine("code-11", StateDefault, ws("    "), ty("src"), ws("  "), ty("string")),
		codeLine("code-12", StateDefault, ws("    "), ty("pos"), ws("  "), ty("int")),
		codeLine("code-13", StateDefault, ws("    "), ty("line"), ws(" "), ty("int")),
		codeLine("code-14", StateDefault, punc("}")),
		codeLine("code-15", StateDefault, com("// NewScanner returns a scanner at the start of src.")),
		codeLine("code-16", StateDefault, kw("func"), ws(" "), fn("NewScanner"), punc("("), ty("src"), ws(" "), ty("string"), punc(")"), ws(" "), punc("*"), ty("Scanner"), ws(" "), punc("{")),
		codeLine("code-17", StateDefault, ws("    "), kw("return"), ws(" "), punc("&"), ty("Scanner"), punc("{"), ty("src"), punc(":"), ws(" "), ty("src"), punc(","), ws(" "), ty("line"), punc(":"), ws(" "), num("1"), punc("}")),
		codeLine("code-18", StateDefault, punc("}")),
		codeLine("code-19", StateDefault, com("// Next advances and returns the next token.")),
		codeLine("code-20", StateDefault, kw("func"), ws(" "), punc("("), ty("s"), ws(" "), punc("*"), ty("Scanner"), punc(")"), ws(" "), fn("Next"), punc("()"), ws(" "), ty("Token").marked(StateSearchMatch).bg(RoleSelectionBackground), ws(" "), punc("{")),
		codeLine("code-21", StateDefault, ws("    "), kw("if"), ws(" "), ty("s"), punc("."), ty("pos"), ws(" "), punc(">="), ws(" "), fn("len"), punc("("), ty("s"), punc("."), ty("src"), punc(")"), ws(" "), punc("{")),
		codeLine("code-22", StateDefault, ws("        "), kw("return"), ws(" "), ty("Token").marked(StateSearchMatch).bg(RoleSelectionBackground), punc("{"), ty("Kind"), punc(":"), ws(" "), str(`"eof"`), punc("}")),
		codeLine("code-23", StateDefault, ws("    "), punc("}")),
		codeLine("code-24", StateActive, ws("    "), ty("ch"), ws(" "), punc(":="), ws(" "), ty("s"), punc("."), ty("src"), punc("["), ty("s"), punc("."), ty("pos"), punc("]")),
		codeLine("code-25", StateDefault, ws("    "), ty("s"), punc("."), ty("pos"), punc("++")),
		codeLine("code-26", StateDefault, ws("    "), kw("if"), ws(" "), ty("ch"), ws(" "), punc("=="), ws(" "), str(`'\n'`), ws(" "), punc("{")),
		codeLine("code-27", StateDefault, ws("        "), ty("s"), punc("."), ty("line"), punc("++")),
		codeLine("code-28", StateDefault, ws("    "), punc("}")),
		codeLine("code-29", StateDefault, ws("    "), kw("return"), ws(" "), ty("Token"), punc("{"), ty("Kind"), punc(":"), ws(" "), str(`"char"`), punc(","), ws(" "), ty("Literal"), punc(":"), ws(" "), fn("string"), punc("("), ty("ch"), punc(")"), punc(","), ws(" "), ty("Line"), punc(":"), ws(" "), ty("s"), punc("."), ty("line"), punc("}")),
		codeLine("code-30", StateDefault, punc("}")),
		codeLine("code-31", StateDefault, com("// Peek returns the next token without advancing.")),
		codeLine("code-32", StateDefault, kw("func"), ws(" "), punc("("), ty("s"), ws(" "), punc("*"), ty("Scanner"), punc(")"), ws(" "), fn("Peek"), punc("()"), ws(" "), ty("Token"), ws(" "), punc("{")),
		codeLine("code-33", StateDefault, ws("    "), ty("saved").marked(StateWarning), ws(" "), punc(":="), ws(" "), ty("s"), punc("."), ty("pos")),
		{ID: "annotation-warning", Kind: ContentText, State: StateWarning, Text: "warning: variable 'saved' is never used", Foreground: RoleWarning},
		codeLine("code-34", StateDefault, ws("    "), ty("token"), ws(" "), punc(":="), ws(" "), ty("s"), punc("."), fn("Next"), punc("()")),
		codeLine("code-35", StateDefault, ws("    "), ty("s"), punc("."), ty("pos"), ws(" "), punc("="), ws(" "), ty("saved")),
		codeLine("code-36", StateDefault, ws("    "), kw("return"), ws(" "), ty("token")),
		codeLine("code-37", StateDefault, punc("}")),
		codeLine("code-38", StateDefault, com("// ScanAll drains the scanner into a slice.")),
		codeLine("code-39", StateDefault, kw("func"), ws(" "), punc("("), ty("s"), ws(" "), punc("*"), ty("Scanner"), punc(")"), ws(" "), fn("ScanAll"), punc("()"), ws(" "), punc("["), punc("]"), ty("Token"), ws(" "), punc("{")),
		codeLine("code-40", StateDefault, ws("    "), kw("var"), ws(" "), ty("out"), ws(" "), punc("["), punc("]"), ty("Token")),
		codeLine("code-41", StateDefault, ws("    "), kw("for"), ws(" "), ty("s"), punc("."), ty("pos"), ws(" "), punc("<"), ws(" "), fn("len"), punc("("), ty("s"), punc("."), ty("src"), punc(")"), ws(" "), punc("{")),
		codeLine("code-42", StateDefault, ws("        "), ty("out"), ws(" "), punc("="), ws(" "), fn("append"), punc("("), ty("out"), punc(","), ws(" "), ty("s"), punc("."), fn("Next"), punc("()"), punc(")")),
		codeLine("code-43", StateDefault, ws("    "), punc("}")),
		codeLine("code-44", StateDefault, ws("    "), kw("return"), ws(" "), ty("out")),
		codeLine("code-45", StateDefault, punc("}")),
		codeLine("code-46", StateDefault, com("// Validate reports the first obvious mistake in the stream.")),
		codeLine("code-47", StateDefault, kw("func"), ws(" "), punc("("), ty("s"), ws(" "), punc("*"), ty("Scanner"), punc(")"), ws(" "), fn("Validate"), punc("()"), ws(" "), ty("error"), ws(" "), punc("{")),
		codeLine("code-48", StateDefault, ws("    "), kw("if"), ws(" "), ty("s"), punc("."), ty("pos"), ws(" "), punc("<"), ws(" "), num("0"), ws(" "), punc("{")),
		codeLine("code-49", StateDefault, ws("        "), kw("return"), ws(" "), ty("fmt"), punc("."), fn("Errorf"), punc("("), str(`"negative position: %d"`), punc(","), ws(" "), ty("s"), punc("."), ty("pos"), punc(")")),
		codeLine("code-50", StateDefault, ws("    "), punc("}")),
		codeLine("code-51", StateDefault, ws("    "), kw("return"), ws(" "), ty("Token"), punc("{"), ty("Kind"), punc(":"), ws(" "), fn("undefinedName").marked(StateError), punc("}")),
		{ID: "annotation-error", Kind: ContentText, State: StateError, Text: "error: undefined name 'undefinedName'", Foreground: RoleError},
		codeLine("code-52", StateDefault, punc("}")),
	}
}

// codeToken is one syntax-highlighted run of a code line.
type codeToken struct {
	text       string
	role       RoleID
	state      SemanticState
	background RoleID
}

func (t codeToken) marked(state SemanticState) codeToken {
	t.state = state
	return t
}

func (t codeToken) bg(role RoleID) codeToken {
	t.background = role
	return t
}

func kw(s string) codeToken  { return codeToken{text: s, role: codeKeywordRole, state: StateDefault} }
func ty(s string) codeToken  { return codeToken{text: s, role: codeTypeRole, state: StateDefault} }
func fn(s string) codeToken  { return codeToken{text: s, role: codeFunctionRole, state: StateDefault} }
func str(s string) codeToken { return codeToken{text: s, role: codeStringRole, state: StateDefault} }
func num(s string) codeToken { return codeToken{text: s, role: codeNumberRole, state: StateDefault} }
func com(s string) codeToken { return codeToken{text: s, role: codeCommentRole, state: StateDefault} }
func punc(s string) codeToken {
	return codeToken{text: s, role: codePunctuationRole, state: StateDefault}
}
func ws(s string) codeToken { return codeToken{text: s, role: RoleForeground, state: StateDefault} }

// codeLine builds a code content block whose runs carry one distinct role per token
// class, which is what exposes syntax density rather than a single flat string.
func codeLine(id string, state SemanticState, tokens ...codeToken) ContentBlock {
	runs := make([]ContentRun, 0, len(tokens))
	for _, token := range tokens {
		runs = append(runs, ContentRun{
			Text:       token.text,
			State:      token.state,
			Background: token.background,
			Foreground: token.role,
		})
	}
	return ContentBlock{ID: ContentBlockID(id), Kind: ContentCode, State: state, Runs: runs}
}

func lineNumberBlocks(count int) []ContentBlock {
	blocks := make([]ContentBlock, 0, count)
	for i := 1; i <= count; i++ {
		blocks = append(blocks, ContentBlock{
			ID:         ContentBlockID("gutter-" + strconv.Itoa(i)),
			Kind:       ContentCode,
			State:      StateDefault,
			Text:       strconv.Itoa(i),
			Foreground: RoleTextMuted,
		})
	}
	return blocks
}
