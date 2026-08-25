package fixture

// ExtremePalette is a development-only palette that is a complete, valid palette input
// while deliberately collapsing one class of visibility or hierarchy distinction. It
// exists so a reviewer can see a failure mode before real candidates do, and it never
// reaches a participant flow.
type ExtremePalette struct {
	ID          PaletteID
	Description string
	Palette     Palette
}

// ExtremePalettes returns the development-only extreme palettes, one per failure class:
// flat hierarchy, collapsed adjacent surfaces, invisible cursor and selection, semantic
// collisions, and ANSI normal/bright collapse. Each retains every required role and
// alias, so each is a valid palette whose failure lives in its values, not its shape.
func ExtremePalettes() []ExtremePalette {
	return []ExtremePalette{
		extremeFlatHierarchy(),
		extremeCollapsedSurfaces(),
		extremeInvisibleCursorSelection(),
		extremeSemanticCollision(),
		extremeANSIBrightCollapse(),
	}
}

// extremeBaseColors is a neutral dark palette that keeps every role distinct. Extreme
// palettes override a few roles on top of it, so the failure is isolated to the roles
// the palette is meant to expose.
func extremeBaseColors() map[RoleID]string {
	return map[RoleID]string{
		RoleBackground:          "#1a1a1a",
		RoleSurface:             "#242424",
		RoleSurfaceElevated:     "#2e2e2e",
		RoleForeground:          "#e0e0e0",
		RoleTextSecondary:       "#a0a0a0",
		RoleTextMuted:           "#6a6a6a",
		RoleAccent:              "#4a9eff",
		RoleAccentForeground:    "#0a0a0a",
		RoleFocus:               "#4a9eff",
		RoleCursor:              "#e0e0e0",
		RoleSelectionBackground: "#3a5a8a",
		RoleSelectionForeground: "#ffffff",
		RoleError:               "#ff5555",
		RoleWarning:             "#ffb86c",
		RoleSuccess:             "#50fa7b",
		RoleInfo:                "#8be9fd",
		RoleDiffAdd:             "#50fa7b",
		RoleDiffRemove:          "#ff5555",
		RoleDiffModify:          "#ffb86c",
		RoleANSI0:               "#1a1a1a",
		RoleANSI1:               "#ff5555",
		RoleANSI2:               "#50fa7b",
		RoleANSI3:               "#ffb86c",
		RoleANSI4:               "#4a9eff",
		RoleANSI5:               "#bd93f9",
		RoleANSI6:               "#8be9fd",
		RoleANSI7:               "#e0e0e0",
		RoleANSI8:               "#6a6a6a",
		RoleANSI9:               "#ff6e6e",
		RoleANSI10:              "#69ff94",
		RoleANSI11:              "#ffd479",
		RoleANSI12:              "#6cb6ff",
		RoleANSI13:              "#d0a9ff",
		RoleANSI14:              "#a4f0ff",
		RoleANSI15:              "#ffffff",
	}
}

func extremeFlatHierarchy() ExtremePalette {
	return extremePalette(
		"extreme-flat-hierarchy",
		"background, surface and surface-elevated collapse to one value, and foreground, text-secondary and text-muted collapse to one value, so hierarchy is invisible.",
		map[RoleID]string{
			RoleSurface:         "#1a1a1a",
			RoleSurfaceElevated: "#1a1a1a",
			RoleTextSecondary:   "#e0e0e0",
			RoleTextMuted:       "#e0e0e0",
		},
	)
}

func extremeCollapsedSurfaces() ExtremePalette {
	return extremePalette(
		"extreme-collapsed-surfaces",
		"surface and surface-elevated share one value, so adjacent surfaces such as a window against a status bar lose their edge.",
		map[RoleID]string{
			RoleSurfaceElevated: "#242424",
		},
	)
}

func extremeInvisibleCursorSelection() ExtremePalette {
	return extremePalette(
		"extreme-invisible-cursor-selection",
		"cursor and selection-background match the background, so the cursor and the selected text are invisible.",
		map[RoleID]string{
			RoleCursor:              "#1a1a1a",
			RoleSelectionBackground: "#1a1a1a",
		},
	)
}

func extremeSemanticCollision() ExtremePalette {
	return extremePalette(
		"extreme-semantic-collision",
		"error and success share one value, and warning and info share one value, so semantic states collide.",
		map[RoleID]string{
			RoleSuccess: "#ff5555",
			RoleInfo:    "#ffb86c",
		},
	)
}

func extremeANSIBrightCollapse() ExtremePalette {
	return extremePalette(
		"extreme-ansi-bright-collapse",
		"each ANSI bright slot shares its normal counterpart, so normal and bright text are indistinguishable.",
		map[RoleID]string{
			RoleANSI8:  "#1a1a1a",
			RoleANSI9:  "#ff5555",
			RoleANSI10: "#50fa7b",
			RoleANSI11: "#ffb86c",
			RoleANSI12: "#4a9eff",
			RoleANSI13: "#bd93f9",
			RoleANSI14: "#8be9fd",
			RoleANSI15: "#e0e0e0",
		},
	)
}

func extremePalette(id PaletteID, description string, overrides map[RoleID]string) ExtremePalette {
	colors := extremeBaseColors()
	for role, color := range overrides {
		colors[role] = color
	}
	return ExtremePalette{
		ID:          id,
		Description: description,
		Palette: Palette{
			SchemaVersion: PaletteSchemaVersion,
			ID:            id,
			SemanticCore:  extremeSemanticCore(colors),
			Terminal: TerminalPalette{
				ANSI:    extremeANSI(colors),
				Aliases: extremeTerminalAliases(),
			},
		},
	}
}

func extremeSemanticCore(colors map[RoleID]string) []Role {
	roles := make([]Role, 0, len(semanticCoreRoleIDs))
	for _, id := range semanticCoreRoleIDs {
		roles = append(roles, Role{ID: id, Value: extremeColor(colors[id])})
	}
	return roles
}

func extremeANSI(colors map[RoleID]string) []Role {
	roles := make([]Role, 0, len(ansiRoleIDs))
	for _, id := range ansiRoleIDs {
		roles = append(roles, Role{ID: id, Value: extremeColor(colors[id])})
	}
	return roles
}

func extremeTerminalAliases() []Role {
	return []Role{
		{ID: RoleTerminalBackground, Alias: &Alias{Target: RoleBackground}},
		{ID: RoleTerminalForeground, Alias: &Alias{Target: RoleForeground}},
		{ID: RoleTerminalCursor, Alias: &Alias{Target: RoleCursor}},
		{ID: RoleTerminalSelectionBackground, Alias: &Alias{Target: RoleSelectionBackground}},
		{ID: RoleTerminalSelectionForeground, Alias: &Alias{Target: RoleSelectionForeground}},
	}
}

// extremeColor authors a role value. The OKLCH is a placeholder: the authoritative
// bounded types and OKLCH-to-sRGB consistency belong to the color package, which does
// not exist yet, and the palette validation only checks finiteness and hex shape.
func extremeColor(srgb string) *AuthoredColor {
	return &AuthoredColor{
		OKLCH: OKLCH{Lightness: 0.5, Chroma: 0, Hue: 0},
		SRGB:  srgb,
		Provenance: Provenance{
			ProfileClaims: []ProfileClaimID{"claim-extreme"},
			Judgments:     []JudgmentID{"judgment-extreme"},
		},
		Validity: ValidityMetadata{
			Evidence: []ValidityEvidenceID{"validity-extreme"},
		},
	}
}
