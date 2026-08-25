// Package fixture defines the versioned measurement surface and source-palette contract.
package fixture

import "encoding/json"

const (
	FixtureSchemaVersion = "1"
	PaletteSchemaVersion = "1"
)

type FixtureSetID string
type SceneID string
type RegionID string
type ContentBlockID string
type PaletteID string
type ProfileClaimID string
type JudgmentID string
type ValidityEvidenceID string

type SceneFamily string

const (
	FamilyTerminalAgent  SceneFamily = "terminal-agent"
	FamilyCodeDiff       SceneFamily = "code-diff"
	FamilyDesktopShell   SceneFamily = "desktop-shell"
	FamilyReadingMonitor SceneFamily = "reading-monitoring"
)

type RegionKind string

const (
	RegionFrame   RegionKind = "frame"
	RegionSurface RegionKind = "surface"
	RegionRow     RegionKind = "row"
	RegionColumn  RegionKind = "column"
	RegionTabs    RegionKind = "tabs"
	RegionList    RegionKind = "list"
	RegionTable   RegionKind = "table"
	RegionStatus  RegionKind = "status"
	RegionOverlay RegionKind = "overlay"
)

type ContentKind string

const (
	ContentText       ContentKind = "text"
	ContentCode       ContentKind = "code"
	ContentTab        ContentKind = "tab"
	ContentListItem   ContentKind = "list-item"
	ContentTableCell  ContentKind = "table-cell"
	ContentStatusItem ContentKind = "status-item"
	ContentSelection  ContentKind = "selection"
	ContentFocusRing  ContentKind = "focus-ring"
	ContentCursor     ContentKind = "cursor"
	ContentBadge      ContentKind = "badge"
)

type SemanticState string

const (
	StateDefault     SemanticState = "default"
	StateActive      SemanticState = "active"
	StateInactive    SemanticState = "inactive"
	StateFocused     SemanticState = "focused"
	StateSelected    SemanticState = "selected"
	StateMuted       SemanticState = "muted"
	StateSuccess     SemanticState = "success"
	StateWarning     SemanticState = "warning"
	StateError       SemanticState = "error"
	StateInfo        SemanticState = "info"
	StateUrgent      SemanticState = "urgent"
	StateAdded       SemanticState = "added"
	StateRemoved     SemanticState = "removed"
	StateModified    SemanticState = "modified"
	StateSearchMatch SemanticState = "search-match"
)

type FixtureSet struct {
	SchemaVersion string       `json:"schemaVersion"`
	ID            FixtureSetID `json:"id"`
	Scenes        []Scene      `json:"scenes"`
}

type Scene struct {
	ID      SceneID     `json:"id"`
	Family  SceneFamily `json:"family"`
	Regions []Region    `json:"regions"`
}

// Region is one ordered structural area. Layout remains a fixed renderer concern, so the
// fixture contract cannot grow an arbitrary tree of user-defined components.
type Region struct {
	ID         RegionID       `json:"id"`
	Kind       RegionKind     `json:"kind"`
	State      SemanticState  `json:"state"`
	Background RoleID         `json:"background,omitempty"`
	Foreground RoleID         `json:"foreground,omitempty"`
	Border     RoleID         `json:"border,omitempty"`
	Blocks     []ContentBlock `json:"blocks"`
}

type ContentBlock struct {
	ID         ContentBlockID `json:"id"`
	Kind       ContentKind    `json:"kind"`
	State      SemanticState  `json:"state"`
	Text       string         `json:"text,omitempty"`
	Runs       []ContentRun   `json:"runs,omitempty"`
	Background RoleID         `json:"background,omitempty"`
	Foreground RoleID         `json:"foreground,omitempty"`
	Border     RoleID         `json:"border,omitempty"`
}

// ContentRun gives a block ordered, role-aware text without making fixture content a
// recursive rendering language.
type ContentRun struct {
	Text       string        `json:"text"`
	State      SemanticState `json:"state"`
	Background RoleID        `json:"background,omitempty"`
	Foreground RoleID        `json:"foreground,omitempty"`
}

// MarshalJSON normalizes a nil collection slice to an empty array before encoding. A
// field declared without omitempty promises the key is always present with a value of its
// declared type, but a nil Go slice marshals to null, which contradicts that promise and
// forces every consumer to defend against a shape the schema said could not occur.
func (f FixtureSet) MarshalJSON() ([]byte, error) {
	type plain FixtureSet
	if f.Scenes == nil {
		f.Scenes = []Scene{}
	}
	return json.Marshal(plain(f))
}

func (s Scene) MarshalJSON() ([]byte, error) {
	type plain Scene
	if s.Regions == nil {
		s.Regions = []Region{}
	}
	return json.Marshal(plain(s))
}

func (r Region) MarshalJSON() ([]byte, error) {
	type plain Region
	if r.Blocks == nil {
		r.Blocks = []ContentBlock{}
	}
	return json.Marshal(plain(r))
}

// MarshalIndent is the canonical human-readable JSON encoding for versioned artifacts.
func MarshalIndent(value any) ([]byte, error) {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(encoded, '\n'), nil
}
