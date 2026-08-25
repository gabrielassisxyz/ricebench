package fixture

import "fmt"

// FixtureSetVersion identifies one frozen version of a fixture set. It is distinct from
// FixtureSchemaVersion, which is a schema compatibility field, not fixture-set identity.
type FixtureSetVersion string

// FixtureSetVersionOne is the first frozen fixture-set version.
const FixtureSetVersionOne FixtureSetVersion = "v1"

// RegistryEntry is the frozen artifact for one fixture-set version. It records the
// fixture-set id and version, the canonical serialization and hash of every scene, the
// canonical hash of the whole set, and the role coverage. It carries no clock-derived
// metadata, so freezing the same content twice produces the same entry.
type RegistryEntry struct {
	ID            FixtureSetID        `json:"id"`
	Version       FixtureSetVersion   `json:"version"`
	SchemaVersion string              `json:"schemaVersion"`
	Hash          string              `json:"hash"`
	Scenes        []RegistryScene     `json:"scenes"`
	Coverage      []RoleCoverageEntry `json:"coverage"`
}

// RegistryScene is the canonical serialization and hash of one scene in a frozen set.
type RegistryScene struct {
	ID        SceneID `json:"id"`
	Canonical string  `json:"canonical"`
	Hash      string  `json:"hash"`
}

// Freeze computes the frozen registry entry for a fixture set at a version. It records
// the fixture-set id and version, the canonical serialization and hash of every scene,
// the canonical hash of the whole set, and the role coverage. It carries no
// clock-derived metadata, so freezing the same content twice produces the same entry.
func Freeze(definition FixtureSet, version FixtureSetVersion) (RegistryEntry, error) {
	if err := definition.Validate(); err != nil {
		return RegistryEntry{}, err
	}
	if err := validateID("version", string(version)); err != nil {
		return RegistryEntry{}, err
	}
	hash, err := CanonicalHash(definition)
	if err != nil {
		return RegistryEntry{}, err
	}
	sceneHashes, err := CanonicalSceneHashes(definition)
	if err != nil {
		return RegistryEntry{}, err
	}
	scenes := make([]RegistryScene, 0, len(definition.Scenes))
	for _, scene := range definition.Scenes {
		canonical, err := CanonicalScene(definition, scene.ID)
		if err != nil {
			return RegistryEntry{}, err
		}
		scenes = append(scenes, RegistryScene{
			ID:        scene.ID,
			Canonical: string(canonical),
			Hash:      sceneHashes[scene.ID],
		})
	}
	return RegistryEntry{
		ID:            definition.ID,
		Version:       version,
		SchemaVersion: definition.SchemaVersion,
		Hash:          hash,
		Scenes:        scenes,
		Coverage:      CoverageView(definition),
	}, nil
}

// Registry holds frozen fixture-set versions, keyed by fixture-set id and version.
type Registry struct {
	entries map[FixtureSetID]map[FixtureSetVersion]RegistryEntry
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{entries: make(map[FixtureSetID]map[FixtureSetVersion]RegistryEntry)}
}

// Freeze records a fixture set as a new version. Freezing the same content under the
// same version is idempotent and returns the existing entry. Freezing different content
// under an existing version is rejected, so changing content requires a new version
// rather than mutating a frozen one.
func (r *Registry) Freeze(definition FixtureSet, version FixtureSetVersion) (RegistryEntry, error) {
	entry, err := Freeze(definition, version)
	if err != nil {
		return RegistryEntry{}, err
	}
	byVersion, ok := r.entries[definition.ID]
	if !ok {
		byVersion = make(map[FixtureSetVersion]RegistryEntry)
		r.entries[definition.ID] = byVersion
	}
	if existing, ok := byVersion[version]; ok {
		if existing.Hash != entry.Hash {
			return RegistryEntry{}, fmt.Errorf("fixture set %q version %q is already frozen with different content", definition.ID, version)
		}
		return existing, nil
	}
	byVersion[version] = entry
	return entry, nil
}

// Entry returns the frozen entry for an id and version, so a later manifest can select
// it without redefining it.
func (r *Registry) Entry(id FixtureSetID, version FixtureSetVersion) (RegistryEntry, bool) {
	byVersion, ok := r.entries[id]
	if !ok {
		return RegistryEntry{}, false
	}
	entry, ok := byVersion[version]
	return entry, ok
}

// FrozenFixtureSets returns a registry with version one of every fixture family frozen.
// The shipped sets are valid and the version is a stable ID, so the freeze cannot fail.
func FrozenFixtureSets() *Registry {
	registry := NewRegistry()
	sets := []FixtureSet{
		TerminalAgentFixtureSet(),
		CodeDiffFixtureSet(),
		DesktopShellFixtureSet(),
		ReadingMonitoringFixture(),
	}
	for _, set := range sets {
		if _, err := registry.Freeze(set, FixtureSetVersionOne); err != nil {
			panic(err)
		}
	}
	return registry
}
