package fixture

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Canonicalize returns the deterministic canonical JSON of a fixture set's structure.
// It derives only from schemaVersion, id, and scenes, retaining role-reference
// identifiers and omitting resolved color values, which live in the palette and never
// enter the fixture contract. A role field that carries a resolved color value instead
// of a role reference is rejected rather than silently canonicalized.
func Canonicalize(definition FixtureSet) ([]byte, error) {
	if err := rejectResolvedRoleValues(definition); err != nil {
		return nil, err
	}
	return json.Marshal(definition)
}

// CanonicalHash returns the hex SHA-256 of the canonical structure.
func CanonicalHash(definition FixtureSet) (string, error) {
	canonical, err := Canonicalize(definition)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:]), nil
}

// CanonicalSceneHashes returns the canonical structure hash for each scene, keyed by
// scene ID. Each hash covers the set's schemaVersion and id plus the scene's full
// structure, so a scene moved between sets changes identity.
func CanonicalSceneHashes(definition FixtureSet) (map[SceneID]string, error) {
	hashes := make(map[SceneID]string, len(definition.Scenes))
	for _, scene := range definition.Scenes {
		hash, err := CanonicalHash(FixtureSet{
			SchemaVersion: definition.SchemaVersion,
			ID:            definition.ID,
			Scenes:        []Scene{scene},
		})
		if err != nil {
			return nil, err
		}
		hashes[scene.ID] = hash
	}
	return hashes, nil
}

// CanonicalScene returns the canonical JSON of a single scene, wrapped in its set
// identity. The representation is indented so a golden stays reviewable.
func CanonicalScene(definition FixtureSet, sceneID SceneID) ([]byte, error) {
	for _, scene := range definition.Scenes {
		if scene.ID != sceneID {
			continue
		}
		wrapped := FixtureSet{
			SchemaVersion: definition.SchemaVersion,
			ID:            definition.ID,
			Scenes:        []Scene{scene},
		}
		if err := rejectResolvedRoleValues(wrapped); err != nil {
			return nil, err
		}
		return MarshalIndent(wrapped)
	}
	return nil, fmt.Errorf("scene %q not found", sceneID)
}

// Equivalent reports whether two fixture sets share the same canonical structure.
func Equivalent(left, right FixtureSet) (bool, error) {
	leftHash, err := CanonicalHash(left)
	if err != nil {
		return false, err
	}
	rightHash, err := CanonicalHash(right)
	if err != nil {
		return false, err
	}
	return leftHash == rightHash, nil
}

// Divergence returns the canonical JSON paths at which two fixture sets differ, such as
// scenes[0].regions[0].background. An empty result means the sets are canonically
// identical.
func Divergence(left, right FixtureSet) ([]string, error) {
	leftCanonical, err := Canonicalize(left)
	if err != nil {
		return nil, err
	}
	rightCanonical, err := Canonicalize(right)
	if err != nil {
		return nil, err
	}
	var leftDoc, rightDoc any
	if err := json.Unmarshal(leftCanonical, &leftDoc); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rightCanonical, &rightDoc); err != nil {
		return nil, err
	}
	var paths []string
	diffPaths("", leftDoc, rightDoc, &paths)
	return paths, nil
}

// diffPaths records the canonical JSON paths at which two decoded documents differ.
func diffPaths(path string, left, right any, out *[]string) {
	switch leftValue := left.(type) {
	case map[string]any:
		rightValue, ok := right.(map[string]any)
		if !ok {
			*out = append(*out, path)
			return
		}
		keys := make([]string, 0, len(leftValue)+len(rightValue))
		seen := make(map[string]struct{}, len(leftValue)+len(rightValue))
		for key := range leftValue {
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				keys = append(keys, key)
			}
		}
		for key := range rightValue {
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				keys = append(keys, key)
			}
		}
		sort.Strings(keys)
		for _, key := range keys {
			childPath := key
			if path != "" {
				childPath = path + "." + key
			}
			leftChild, leftOK := leftValue[key]
			rightChild, rightOK := rightValue[key]
			switch {
			case !leftOK || !rightOK:
				*out = append(*out, childPath)
			default:
				diffPaths(childPath, leftChild, rightChild, out)
			}
		}
	case []any:
		rightValue, ok := right.([]any)
		if !ok {
			*out = append(*out, path)
			return
		}
		length := len(leftValue)
		if len(rightValue) > length {
			length = len(rightValue)
		}
		for index := 0; index < length; index++ {
			childPath := fmt.Sprintf("%s[%d]", path, index)
			switch {
			case index >= len(leftValue) || index >= len(rightValue):
				*out = append(*out, childPath)
			default:
				diffPaths(childPath, leftValue[index], rightValue[index], out)
			}
		}
	default:
		if left != right {
			*out = append(*out, path)
		}
	}
}

// rejectResolvedRoleValues rejects a fixture whose role-reference fields carry a
// resolved color value instead of a role reference. A resolved value bypasses the
// role-reference mechanism: the renderer would treat it as a literal color, and the
// canonical structure would then depend on a palette-specific value rather than a role
// identity.
func rejectResolvedRoleValues(definition FixtureSet) error {
	for sceneIndex, scene := range definition.Scenes {
		for regionIndex, region := range scene.Regions {
			path := fmt.Sprintf("scenes[%d].regions[%d]", sceneIndex, regionIndex)
			if err := rejectResolvedRole(path+".background", region.Background); err != nil {
				return err
			}
			if err := rejectResolvedRole(path+".foreground", region.Foreground); err != nil {
				return err
			}
			if err := rejectResolvedRole(path+".border", region.Border); err != nil {
				return err
			}
			for blockIndex, block := range region.Blocks {
				blockPath := fmt.Sprintf("%s.blocks[%d]", path, blockIndex)
				if err := rejectResolvedRole(blockPath+".background", block.Background); err != nil {
					return err
				}
				if err := rejectResolvedRole(blockPath+".foreground", block.Foreground); err != nil {
					return err
				}
				if err := rejectResolvedRole(blockPath+".border", block.Border); err != nil {
					return err
				}
				for runIndex, run := range block.Runs {
					runPath := fmt.Sprintf("%s.runs[%d]", blockPath, runIndex)
					if err := rejectResolvedRole(runPath+".background", run.Background); err != nil {
						return err
					}
					if err := rejectResolvedRole(runPath+".foreground", run.Foreground); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func rejectResolvedRole(path string, role RoleID) error {
	if role == "" {
		return nil
	}
	if isResolvedColorValue(role) {
		return fmt.Errorf("%s: resolved color value %q bypasses role reference", path, role)
	}
	return nil
}

// isResolvedColorValue reports whether a role field carries a literal color rather than
// a role reference. The fixture contract has no field for a resolved color, so a hex
// color or an OKLCH literal in a role field is a bypass, not a reference.
func isResolvedColorValue(role RoleID) bool {
	value := string(role)
	if hexColorPattern.MatchString(value) {
		return true
	}
	lower := strings.ToLower(value)
	return strings.HasPrefix(lower, "oklch(") || strings.HasPrefix(lower, "rgb(") || strings.HasPrefix(lower, "hsl(")
}
