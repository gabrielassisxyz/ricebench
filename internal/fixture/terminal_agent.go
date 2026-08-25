package fixture

// TerminalAgentFixtureSet returns the terminal-agent scene family: a shell session, an
// agent conversation, and an ANSI role exercise. Every scene renders from this one
// structural definition, so a candidate comparison can only differ in color.
func TerminalAgentFixtureSet() FixtureSet {
	return FixtureSet{
		SchemaVersion: FixtureSchemaVersion,
		ID:            "fixture-terminal-agent",
		Scenes: []Scene{
			terminalShellScene(),
			agentConversationScene(),
			terminalANSIScene(),
		},
	}
}

// terminalShellScene exercises prompt, path, command, normal output and the four
// semantic output states, plus cursor, selection, active and inactive content.
func terminalShellScene() Scene {
	return Scene{
		ID:     "terminal-shell",
		Family: FamilyTerminalAgent,
		Regions: []Region{
			{
				ID:         "shell-frame",
				Kind:       RegionFrame,
				State:      StateActive,
				Background: RoleTerminalBackground,
				Foreground: RoleTerminalForeground,
				Border:     RoleSurface,
				Blocks: []ContentBlock{
					{
						ID:    "shell-prompt",
						Kind:  ContentText,
						State: StateDefault,
						Runs: []ContentRun{
							{Text: "❯ ", State: StateDefault, Foreground: RoleANSI2},
							{Text: "/srv/app", State: StateDefault, Foreground: RoleANSI4},
							{Text: " (main)", State: StateMuted, Foreground: RoleTextMuted},
							{Text: " make test", State: StateDefault, Foreground: RoleTerminalForeground, Background: RoleTerminalBackground},
						},
					},
					{
						ID:         "shell-output",
						Kind:       ContentText,
						State:      StateDefault,
						Text:       "running 42 checks across 3 packages",
						Foreground: RoleTerminalForeground,
					},
					{
						ID:         "shell-success",
						Kind:       ContentText,
						State:      StateSuccess,
						Text:       "ok      pkg/parser    0.412s",
						Foreground: RoleSuccess,
					},
					{
						ID:         "shell-warning",
						Kind:       ContentText,
						State:      StateWarning,
						Text:       "warning: flag -coverprofile is deprecated and will be removed",
						Foreground: RoleWarning,
					},
					{
						ID:         "shell-error",
						Kind:       ContentText,
						State:      StateError,
						Text:       "FAIL: expected 3 assertions, got 2",
						Foreground: RoleError,
					},
					{
						ID:         "shell-info",
						Kind:       ContentText,
						State:      StateInfo,
						Text:       "note: 1 test skipped because it needs a display",
						Foreground: RoleInfo,
					},
					{
						ID:    "shell-cursor",
						Kind:  ContentCursor,
						State: StateFocused,
						Runs: []ContentRun{
							{Text: "▊", State: StateFocused, Foreground: RoleTerminalCursor},
						},
					},
					{
						ID:    "shell-selection",
						Kind:  ContentSelection,
						State: StateSelected,
						Runs: []ContentRun{
							{Text: "selected output line", State: StateSelected, Foreground: RoleTerminalSelectionForeground, Background: RoleTerminalSelectionBackground},
						},
					},
					{
						ID:         "shell-active-line",
						Kind:       ContentText,
						State:      StateActive,
						Text:       "active: the command currently running",
						Foreground: RoleTerminalForeground,
					},
					{
						ID:         "shell-inactive",
						Kind:       ContentText,
						State:      StateInactive,
						Text:       "inactive: output from an earlier session, dimmed",
						Foreground: RoleTextMuted,
					},
				},
			},
		},
	}
}

// agentConversationScene exercises long prose, inline and fenced code, a tool-call
// boundary with muted metadata, and active against inactive content.
func agentConversationScene() Scene {
	return Scene{
		ID:     "agent-conversation",
		Family: FamilyTerminalAgent,
		Regions: []Region{
			{
				ID:         "conversation-frame",
				Kind:       RegionFrame,
				State:      StateActive,
				Background: RoleTerminalBackground,
				Foreground: RoleTerminalForeground,
				Border:     RoleSurface,
				Blocks: []ContentBlock{
					{
						ID:    "assistant-prose",
						Kind:  ContentText,
						State: StateDefault,
						Runs: []ContentRun{
							{Text: "The fixture set is a measurement contract. Each scene renders from one structural definition, so a comparison can only differ in color. ", State: StateDefault, Foreground: RoleTerminalForeground},
							{Text: "`FixtureSet.Validate`", State: StateDefault, Foreground: RoleANSI6},
							{Text: " rejects a scene whose family is unknown, and ", State: StateDefault, Foreground: RoleTerminalForeground},
							{Text: "`ValidatePair`", State: StateDefault, Foreground: RoleANSI6},
							{Text: " proves every role reference resolves against a palette.", State: StateDefault, Foreground: RoleTerminalForeground},
						},
					},
					{
						ID:         "assistant-code",
						Kind:       ContentCode,
						State:      StateDefault,
						Background: RoleSurfaceElevated,
						Foreground: RoleTerminalForeground,
						Runs: []ContentRun{
							{Text: "func ", State: StateDefault, Foreground: RoleANSI5},
							{Text: "Validate", State: StateDefault, Foreground: RoleANSI4},
							{Text: "() error {", State: StateDefault, Foreground: RoleTerminalForeground},
							{Text: "  return nil", State: StateDefault, Foreground: RoleTerminalForeground},
							{Text: "}", State: StateDefault, Foreground: RoleTerminalForeground},
						},
					},
					{
						ID:    "tool-call",
						Kind:  ContentText,
						State: StateInfo,
						Runs: []ContentRun{
							{Text: "▸ tool call: ", State: StateInfo, Foreground: RoleInfo},
							{Text: "read_file", State: StateDefault, Foreground: RoleANSI6},
							{Text: " · ", State: StateMuted, Foreground: RoleTextMuted},
							{Text: "path=pkg/parser/schema.go", State: StateMuted, Foreground: RoleTextMuted},
						},
					},
					{
						ID:         "tool-result",
						Kind:       ContentText,
						State:      StateMuted,
						Text:       "result: 214 lines, 6.1 kB, cached",
						Foreground: RoleTextMuted,
					},
					{
						ID:         "conversation-active",
						Kind:       ContentText,
						State:      StateActive,
						Text:       "active: the response currently streaming",
						Foreground: RoleTerminalForeground,
					},
					{
						ID:         "conversation-inactive",
						Kind:       ContentText,
						State:      StateInactive,
						Text:       "inactive: an earlier exchange, dimmed",
						Foreground: RoleTextMuted,
					},
				},
			},
		},
	}
}

// terminalANSIScene embeds the sixteen ANSI roles in a colored status log rather than a
// swatch grid, and places ANSI 0 text on an ANSI background to expose the collision the
// family exists to catch.
func terminalANSIScene() Scene {
	return Scene{
		ID:     "terminal-ansi",
		Family: FamilyTerminalAgent,
		Regions: []Region{
			{
				ID:         "ansi-frame",
				Kind:       RegionFrame,
				State:      StateActive,
				Background: RoleTerminalBackground,
				Foreground: RoleTerminalForeground,
				Border:     RoleSurface,
				Blocks: []ContentBlock{
					{
						ID:    "ansi-log",
						Kind:  ContentCode,
						State: StateDefault,
						Runs: []ContentRun{
							{Text: "PASS ", State: StateSuccess, Foreground: RoleANSI2},
							{Text: "pkg/parser", State: StateDefault, Foreground: RoleANSI7},
							{Text: "  0.412s", State: StateMuted, Foreground: RoleANSI8},
							{Text: "FAIL ", State: StateError, Foreground: RoleANSI1},
							{Text: "pkg/color", State: StateDefault, Foreground: RoleANSI7},
							{Text: "  0.081s", State: StateMuted, Foreground: RoleANSI8},
							{Text: "SKIP ", State: StateWarning, Foreground: RoleANSI3},
							{Text: "pkg/server", State: StateDefault, Foreground: RoleANSI7},
							{Text: "  no display", State: StateMuted, Foreground: RoleANSI8},
							{Text: "INFO ", State: StateInfo, Foreground: RoleANSI4},
							{Text: "listening on :8080", State: StateDefault, Foreground: RoleANSI7},
							{Text: "WARN ", State: StateWarning, Foreground: RoleANSI5},
							{Text: "deprecated flag", State: StateDefault, Foreground: RoleANSI7},
							{Text: "LINK ", State: StateInfo, Foreground: RoleANSI6},
							{Text: "http://localhost:8080", State: StateDefault, Foreground: RoleANSI7},
							{Text: "ERR! ", State: StateError, Foreground: RoleANSI9},
							{Text: "bright error", State: StateDefault, Foreground: RoleANSI15},
							{Text: "OK   ", State: StateSuccess, Foreground: RoleANSI10},
							{Text: "bright success", State: StateDefault, Foreground: RoleANSI15},
							{Text: "WARN!", State: StateWarning, Foreground: RoleANSI11},
							{Text: "bright warning", State: StateDefault, Foreground: RoleANSI15},
							{Text: "INFO!", State: StateInfo, Foreground: RoleANSI12},
							{Text: "bright info", State: StateDefault, Foreground: RoleANSI15},
							{Text: "HILITE ", State: StateDefault, Foreground: RoleANSI13},
							{Text: "bright highlight", State: StateDefault, Foreground: RoleANSI15},
							{Text: "LINK!", State: StateInfo, Foreground: RoleANSI14},
							{Text: "bright link", State: StateDefault, Foreground: RoleANSI15},
						},
					},
					{
						ID:    "ansi-background",
						Kind:  ContentText,
						State: StateDefault,
						Runs: []ContentRun{
							{Text: "normal ", State: StateDefault, Foreground: RoleANSI7},
							{Text: "on blue", State: StateDefault, Foreground: RoleANSI0, Background: RoleANSI4},
							{Text: " after", State: StateDefault, Foreground: RoleANSI7},
						},
					},
				},
			},
		},
	}
}
