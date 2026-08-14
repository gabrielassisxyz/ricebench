package fixture

type RoleID string

const (
	RoleBackground          RoleID = "background"
	RoleSurface             RoleID = "surface"
	RoleSurfaceElevated     RoleID = "surface-elevated"
	RoleForeground          RoleID = "foreground"
	RoleTextSecondary       RoleID = "text-secondary"
	RoleTextMuted           RoleID = "text-muted"
	RoleAccent              RoleID = "accent"
	RoleAccentForeground    RoleID = "accent-foreground"
	RoleFocus               RoleID = "focus"
	RoleCursor              RoleID = "cursor"
	RoleSelectionBackground RoleID = "selection-background"
	RoleSelectionForeground RoleID = "selection-foreground"
	RoleError               RoleID = "error"
	RoleWarning             RoleID = "warning"
	RoleSuccess             RoleID = "success"
	RoleInfo                RoleID = "info"
	RoleDiffAdd             RoleID = "diff-add"
	RoleDiffRemove          RoleID = "diff-remove"
	RoleDiffModify          RoleID = "diff-modify"

	RoleANSI0  RoleID = "ansi-0"
	RoleANSI1  RoleID = "ansi-1"
	RoleANSI2  RoleID = "ansi-2"
	RoleANSI3  RoleID = "ansi-3"
	RoleANSI4  RoleID = "ansi-4"
	RoleANSI5  RoleID = "ansi-5"
	RoleANSI6  RoleID = "ansi-6"
	RoleANSI7  RoleID = "ansi-7"
	RoleANSI8  RoleID = "ansi-8"
	RoleANSI9  RoleID = "ansi-9"
	RoleANSI10 RoleID = "ansi-10"
	RoleANSI11 RoleID = "ansi-11"
	RoleANSI12 RoleID = "ansi-12"
	RoleANSI13 RoleID = "ansi-13"
	RoleANSI14 RoleID = "ansi-14"
	RoleANSI15 RoleID = "ansi-15"

	RoleTerminalBackground          RoleID = "terminal-background"
	RoleTerminalForeground          RoleID = "terminal-foreground"
	RoleTerminalCursor              RoleID = "terminal-cursor"
	RoleTerminalSelectionBackground RoleID = "terminal-selection-background"
	RoleTerminalSelectionForeground RoleID = "terminal-selection-foreground"
)

var semanticCoreRoleIDs = []RoleID{
	RoleBackground, RoleSurface, RoleSurfaceElevated,
	RoleForeground, RoleTextSecondary, RoleTextMuted,
	RoleAccent, RoleAccentForeground, RoleFocus, RoleCursor,
	RoleSelectionBackground, RoleSelectionForeground,
	RoleError, RoleWarning, RoleSuccess, RoleInfo,
	RoleDiffAdd, RoleDiffRemove, RoleDiffModify,
}

var ansiRoleIDs = []RoleID{
	RoleANSI0, RoleANSI1, RoleANSI2, RoleANSI3,
	RoleANSI4, RoleANSI5, RoleANSI6, RoleANSI7,
	RoleANSI8, RoleANSI9, RoleANSI10, RoleANSI11,
	RoleANSI12, RoleANSI13, RoleANSI14, RoleANSI15,
}

var terminalAliasRoleIDs = []RoleID{
	RoleTerminalBackground,
	RoleTerminalForeground,
	RoleTerminalCursor,
	RoleTerminalSelectionBackground,
	RoleTerminalSelectionForeground,
}

type Palette struct {
	SchemaVersion string          `json:"schemaVersion"`
	ID            PaletteID       `json:"id"`
	SemanticCore  []Role          `json:"semanticCore"`
	Terminal      TerminalPalette `json:"terminal"`
}

type TerminalPalette struct {
	ANSI    []Role `json:"ansi"`
	Aliases []Role `json:"aliases"`
}

// Role contains an authored value or an explicit alias. Validate rejects both and neither.
type Role struct {
	ID    RoleID         `json:"id"`
	Value *AuthoredColor `json:"value,omitempty"`
	Alias *Alias         `json:"alias,omitempty"`
}

type Alias struct {
	Target RoleID `json:"target"`
}

type AuthoredColor struct {
	OKLCH      OKLCH            `json:"oklch"`
	SRGB       string           `json:"srgb"`
	Provenance Provenance       `json:"provenance"`
	Validity   ValidityMetadata `json:"validity"`
}

type OKLCH struct {
	Lightness float64 `json:"lightness"`
	Chroma    float64 `json:"chroma"`
	Hue       float64 `json:"hue"`
}

type Provenance struct {
	// Candidate palettes exist before judgments. Final source-palette export is responsible
	// for requiring and validating complete evidence links.
	ProfileClaims []ProfileClaimID `json:"profileClaims"`
	Judgments     []JudgmentID     `json:"judgments"`
}

// ValidityMetadata links a role to the versioned reports that evaluated it. The report
// shape and numeric thresholds belong to the later validity package, not this schema.
type ValidityMetadata struct {
	Evidence []ValidityEvidenceID `json:"evidence"`
}
