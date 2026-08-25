package fixture

import "fmt"

// ReadingMonitoringFixture returns the reading-monitoring scene family as a complete
// fixture set: a long-form reading surface with note chrome, and a monitoring dashboard
// with a dense metric table and a process monitor. The family is the only one that
// exercises duration: glare, muted text that falls below usable, and a hierarchy that
// collapses over a long page only show up in a scene too long to take in at a glance.
func ReadingMonitoringFixture() FixtureSet {
	return FixtureSet{
		SchemaVersion: FixtureSchemaVersion,
		ID:            "fixture-reading-monitoring",
		Scenes: []Scene{
			readingNoteScene(),
			monitoringDashboardScene(),
		},
	}
}

func readingNoteScene() Scene {
	return Scene{
		ID:     "reading-note",
		Family: FamilyReadingMonitor,
		Regions: []Region{
			noteChromeRegion(),
			proseRegion(),
		},
	}
}

func monitoringDashboardScene() Scene {
	regions := []Region{
		processMonitorRegion(),
		metricTableRegion(),
		metricHeaderRegion(),
	}
	regions = append(regions, metricRowRegions()...)
	regions = append(regions, metricScrollbarRegion())
	return Scene{
		ID:      "monitoring-dashboard",
		Family:  FamilyReadingMonitor,
		Regions: regions,
	}
}

// noteChromeRegion models the frame around a reading surface: one active tab and two
// inactive ones. The states are the point, not any particular consumer's chrome.
func noteChromeRegion() Region {
	return Region{
		ID:         "note-chrome",
		Kind:       RegionTabs,
		State:      StateActive,
		Background: RoleSurface,
		Foreground: RoleForeground,
		Blocks: []ContentBlock{
			{ID: "tab-reading", Kind: ContentTab, State: StateActive, Text: "Reading", Foreground: RoleAccent},
			{ID: "tab-notes", Kind: ContentTab, State: StateInactive, Text: "Notes", Foreground: RoleTextMuted},
			{ID: "tab-archive", Kind: ContentTab, State: StateInactive, Text: "Archive", Foreground: RoleTextMuted},
		},
	}
}

// proseRegion carries the long-form reading surface. Body lines are stored one per block
// so the fixture contract, not the renderer, owns the line length: each body line is
// between 60 and 90 characters, and the body carries more than 400 words, so sustained
// reading is judgeable rather than sampled at thumbnail scale.
func proseRegion() Region {
	return Region{
		ID:         "prose",
		Kind:       RegionColumn,
		State:      StateActive,
		Background: RoleSurface,
		Foreground: RoleForeground,
		Blocks:     proseBlocks(),
	}
}

func proseBlocks() []ContentBlock {
	return []ContentBlock{
		textBlock("heading-title", StateActive, RoleAccent, "Keeping a personal archive that lasts"),
		textBlock("meta-published", StateMuted, RoleTextMuted, "Published 2026-08-13 · 2 min read"),
		textBlock("heading-why", StateActive, RoleAccent, "Why plain text wins"),

		textBlock("para-01", StateDefault, RoleForeground, "A personal archive is only as durable as the format it is stored in."),
		runBlock("para-02", StateDefault,
			run("Proprietary formats come and go, but ", StateDefault, RoleForeground),
			run("plain text", StateInfo, RoleAccent),
			run(" has outlived every", StateDefault, RoleForeground),
		),
		textBlock("para-03", StateDefault, RoleForeground, "editor, every operating system, and every vendor that promised to"),
		textBlock("para-04", StateDefault, RoleForeground, "keep your notes safe. When a tool dies, its files still open in anything."),

		runBlock("para-05", StateDefault,
			run("Plain text also makes the archive searchable. Tools like ", StateDefault, RoleForeground),
			run("grep", StateDefault, RoleTextSecondary),
			run(" and ", StateDefault, RoleForeground),
			run("ripgrep", StateDefault, RoleTextSecondary),
			run(" can scan a", StateDefault, RoleForeground),
		),
		textBlock("para-06", StateDefault, RoleForeground, "decade of notes in under a second, and the results are reproducible. A binary format"),
		textBlock("para-07", StateDefault, RoleForeground, "hides its contents behind an application; a text file shows exactly what it holds."),

		textBlock("para-08", StateDefault, RoleForeground, "The trade-off is that plain text has no styling. Headings, links,"),
		textBlock("para-09", StateDefault, RoleForeground, "and emphasis are conventions rather than features. That is a"),
		textBlock("para-10", StateDefault, RoleForeground, "feature in disguise: the meaning lives in the words, not in the"),
		textBlock("para-11", StateDefault, RoleForeground, "rendering, so the archive survives any change of taste in how it is displayed."),

		textBlock("heading-example", StateActive, RoleAccent, "A minimal example"),

		textBlock("para-12", StateDefault, RoleForeground, "A single directory with a few conventions is enough to start. Keep one file"),
		textBlock("para-13", StateDefault, RoleForeground, "per topic, name it with a date prefix, and write in short paragraphs. The"),
		textBlock("para-14", StateDefault, RoleForeground, "structure below has served well for years and needs no special tooling to read."),

		{
			ID:         "code-example",
			Kind:       ContentCode,
			State:      StateDefault,
			Text:       "archive/\n  2026-08-13-reading.md\n  2026-08-14-backups.md\n  index.md",
			Background: RoleSurfaceElevated,
			Foreground: RoleForeground,
		},

		runBlock("para-15", StateDefault,
			run("The ", StateDefault, RoleForeground),
			run("index.md", StateDefault, RoleTextSecondary),
			run(" file is the only piece of structure that matters. It lists", StateDefault, RoleForeground),
		),
		textBlock("para-16", StateDefault, RoleForeground, "every note with a one-line summary, so the archive can be navigated"),
		textBlock("para-17", StateDefault, RoleForeground, "without opening each file. Regenerate it with a script when it drifts out of date."),

		runBlock("para-18", StateDefault,
			run("Versioning is the second habit. A ", StateDefault, RoleForeground),
			run("git", StateDefault, RoleTextSecondary),
			run(" repository turns every edit", StateDefault, RoleForeground),
		),
		textBlock("para-19", StateDefault, RoleForeground, "into a recoverable point, and a remote copy turns the archive into a"),
		textBlock("para-20", StateDefault, RoleForeground, "backup. Neither step requires the notes to be in any particular format."),

		textBlock("para-21", StateDefault, RoleForeground, "A backup that is never tested is not a backup. Restore a single file to a clean"),
		textBlock("para-22", StateDefault, RoleForeground, "directory once a season, and confirm the index still opens. The exercise takes"),
		textBlock("para-23", StateDefault, RoleForeground, "minutes and catches the silent failures that only show up when the archive is needed most."),

		textBlock("heading-breaks", StateActive, RoleAccent, "What breaks an archive"),

		textBlock("para-24", StateDefault, RoleForeground, "The archive fails in predictable ways. Files are lost to a disk that dies, a"),
		textBlock("para-25", StateDefault, RoleForeground, "sync that overwrites, or a rename that breaks a link. Each failure has a"),
		textBlock("para-26", StateDefault, RoleForeground, "cheap fix, but only if the archive is plain text and versioned from the start."),

		textBlock("para-27", StateDefault, RoleForeground, "Links between notes are the most fragile part. A link that names a file"),
		textBlock("para-28", StateDefault, RoleForeground, "by its full path breaks the moment the file moves. Prefer links that"),
		runBlock("para-29", StateDefault,
			run("name ", StateDefault, RoleForeground),
			run("a stable identifier", StateInfo, RoleAccent),
			run(", and let a small script resolve them at read time.", StateDefault, RoleForeground),
		),

		textBlock("para-30", StateDefault, RoleForeground, "The lesson is the same one the archive teaches about itself:"),
		textBlock("para-31", StateDefault, RoleForeground, "keep the format boring, keep the structure explicit, and keep"),
		textBlock("para-32", StateDefault, RoleForeground, "a copy somewhere else. Everything else is a matter of taste."),

		textBlock("meta-footer", StateMuted, RoleTextMuted, "Last updated 2026-08-13 · source: local notes"),
	}
}

// processMonitorRegion carries the densest concentration of competing semantic colors in
// a real desktop. The six states sit side by side so a palette whose warning and error
// roles are too close, or whose informational role reads as normal, fails here first.
func processMonitorRegion() Region {
	return Region{
		ID:         "process-monitor",
		Kind:       RegionStatus,
		State:      StateDefault,
		Background: RoleSurface,
		Foreground: RoleForeground,
		Blocks: []ContentBlock{
			{ID: "proc-normal", Kind: ContentStatusItem, State: StateDefault, Text: "idle", Foreground: RoleForeground},
			{ID: "proc-warning", Kind: ContentStatusItem, State: StateWarning, Text: "load 4.2", Foreground: RoleWarning},
			{ID: "proc-error", Kind: ContentStatusItem, State: StateError, Text: "disk full", Foreground: RoleError},
			{ID: "proc-success", Kind: ContentStatusItem, State: StateSuccess, Text: "backup ok", Foreground: RoleSuccess},
			{ID: "proc-muted", Kind: ContentStatusItem, State: StateMuted, Text: "uptime 12d", Foreground: RoleTextMuted},
			{ID: "proc-info", Kind: ContentStatusItem, State: StateInfo, Text: "sync 3m ago", Foreground: RoleInfo},
		},
	}
}

func metricTableRegion() Region {
	return Region{
		ID:         "metric-table",
		Kind:       RegionTable,
		State:      StateDefault,
		Background: RoleSurface,
		Foreground: RoleForeground,
		Blocks: []ContentBlock{
			textBlock("metric-title", StateActive, RoleAccent, "Service metrics"),
		},
	}
}

func metricHeaderRegion() Region {
	return Region{
		ID:         "metric-header",
		Kind:       RegionRow,
		State:      StateActive,
		Background: RoleSurfaceElevated,
		Foreground: RoleForeground,
		Blocks: []ContentBlock{
			{ID: "metric-header-name", Kind: ContentTableCell, State: StateActive, Text: "Service", Foreground: RoleForeground},
			{ID: "metric-header-cpu", Kind: ContentTableCell, State: StateActive, Text: "CPU", Foreground: RoleForeground},
			{ID: "metric-header-mem", Kind: ContentTableCell, State: StateActive, Text: "Memory", Foreground: RoleForeground},
		},
	}
}

// metricRowRegions returns twenty data rows. Row seven is selected, and the scrollbar
// region that follows marks the list as scrolled rather than showing its top.
func metricRowRegions() []Region {
	services := []struct {
		name string
		cpu  string
		mem  string
	}{
		{"api", "12%", "180M"},
		{"db", "34%", "1.2G"},
		{"cache", "8%", "96M"},
		{"queue", "21%", "240M"},
		{"search", "45%", "890M"},
		{"auth", "6%", "64M"},
		{"web", "18%", "210M"},
		{"worker", "52%", "1.5G"},
		{"mail", "9%", "88M"},
		{"log", "15%", "160M"},
		{"cron", "3%", "32M"},
		{"proxy", "27%", "310M"},
		{"media", "39%", "760M"},
		{"index", "22%", "280M"},
		{"sync", "11%", "120M"},
		{"report", "30%", "420M"},
		{"notify", "7%", "72M"},
		{"backup", "41%", "640M"},
		{"stats", "16%", "190M"},
		{"gate", "5%", "48M"},
	}

	rows := make([]Region, 0, len(services))
	for index, service := range services {
		state := StateDefault
		background := RoleID("")
		foreground := RoleForeground
		if index == 6 {
			state = StateSelected
			background = RoleSelectionBackground
			foreground = RoleSelectionForeground
		}
		rows = append(rows, Region{
			ID:         RegionID(metricRowID(index)),
			Kind:       RegionRow,
			State:      state,
			Background: background,
			Foreground: foreground,
			Blocks: []ContentBlock{
				{ID: ContentBlockID(metricRowID(index) + "-name"), Kind: ContentTableCell, State: state, Text: service.name, Foreground: foreground},
				{ID: ContentBlockID(metricRowID(index) + "-cpu"), Kind: ContentTableCell, State: state, Text: service.cpu, Foreground: foreground},
				{ID: ContentBlockID(metricRowID(index) + "-mem"), Kind: ContentTableCell, State: state, Text: service.mem, Foreground: foreground},
			},
		})
	}
	return rows
}

func metricScrollbarRegion() Region {
	return Region{
		ID:         "metric-scrollbar",
		Kind:       RegionOverlay,
		State:      StateDefault,
		Background: RoleSurfaceElevated,
		Foreground: RoleTextMuted,
		Blocks: []ContentBlock{
			{ID: "metric-scrollbar-thumb", Kind: ContentSelection, State: StateDefault, Text: "scroll", Background: RoleTextMuted},
		},
	}
}

func metricRowID(index int) string {
	return fmt.Sprintf("metric-row-%02d", index+1)
}

func textBlock(id string, state SemanticState, foreground RoleID, text string) ContentBlock {
	return ContentBlock{
		ID:         ContentBlockID(id),
		Kind:       ContentText,
		State:      state,
		Text:       text,
		Foreground: foreground,
	}
}

func runBlock(id string, state SemanticState, runs ...ContentRun) ContentBlock {
	return ContentBlock{
		ID:    ContentBlockID(id),
		Kind:  ContentText,
		State: state,
		Runs:  runs,
	}
}

func run(text string, state SemanticState, foreground RoleID) ContentRun {
	return ContentRun{Text: text, State: state, Foreground: foreground}
}
