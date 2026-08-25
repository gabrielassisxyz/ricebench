// Package store owns the experiment data root and the path boundary around it.
//
// Every experiment operation receives an explicit root and resolves paths through this
// package, so nothing reads or writes outside it. There is deliberately no implicit
// home, current-directory or host fallback: the root is the one fact an operator has to
// state, and a missing or invalid root is an error, not a guess.
package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrRootOmitted reports an empty data root. The operation needs an explicit
	// --data-dir and must not fall back to a home, current-directory or host path.
	ErrRootOmitted = errors.New("data root omitted")

	// ErrRootMissing reports a data root that does not exist.
	ErrRootMissing = errors.New("data root does not exist")

	// ErrRootInaccessible reports a data root that exists but is not an accessible
	// directory.
	ErrRootInaccessible = errors.New("data root is not an accessible directory")

	// ErrPathInvalid reports a path that is not a valid relative name.
	ErrPathInvalid = errors.New("invalid experiment path")

	// ErrPathEscapes reports a path that resolves outside the data root.
	ErrPathEscapes = errors.New("path escapes the data root")
)

// Kind names one family of path the experiment stores beneath the root.
type Kind string

const (
	KindExperiment Kind = "experiment"
	KindManifest   Kind = "manifest"
	KindEvent      Kind = "event"
	KindDerived    Kind = "derived"
	KindPalette    Kind = "palette"
)

// Root is an explicit, validated experiment data directory. Every path the experiment
// reads or writes is resolved against it and cannot escape it.
type Root struct {
	dir string
}

// Resolve validates dir as an explicit data root and returns it. An omitted, missing
// or inaccessible root is an error; there is no fallback and no implicit creation.
func Resolve(dir string) (Root, error) {
	if dir == "" {
		return Root{}, ErrRootOmitted
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return Root{}, fmt.Errorf("resolve data root %q: %w", dir, err)
	}
	abs = filepath.Clean(abs)
	info, err := os.Stat(abs)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Root{}, fmt.Errorf("%w: %q", ErrRootMissing, dir)
		}
		return Root{}, fmt.Errorf("%w: %q: %v", ErrRootInaccessible, dir, err)
	}
	if !info.IsDir() {
		return Root{}, fmt.Errorf("%w: %q is not a directory", ErrRootInaccessible, dir)
	}
	return Root{dir: abs}, nil
}

// Dir returns the absolute, cleaned root path.
func (r Root) Dir() string {
	return r.dir
}

// Path resolves a relative name of the given kind to an absolute path beneath the
// root. An empty, absolute or NUL-containing name, an unknown kind or a traversal is a
// structured error.
func (r Root) Path(kind Kind, name string) (string, error) {
	if !knownKind(kind) {
		return "", fmt.Errorf("%w: unknown kind %q", ErrPathInvalid, kind)
	}
	if name == "" {
		return "", fmt.Errorf("%w: %s name is empty", ErrPathInvalid, kind)
	}
	if strings.ContainsRune(name, '\x00') {
		return "", fmt.Errorf("%w: %s name contains a NUL byte", ErrPathInvalid, kind)
	}
	if filepath.IsAbs(name) {
		return "", fmt.Errorf("%w: %s name %q is absolute", ErrPathInvalid, kind, name)
	}
	clean := filepath.Clean(name)
	if clean == "." {
		return "", fmt.Errorf("%w: %s name %q is empty", ErrPathInvalid, kind, name)
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: %s path %q escapes the data root", ErrPathEscapes, kind, name)
	}
	return filepath.Join(r.dir, clean), nil
}

func knownKind(kind Kind) bool {
	switch kind {
	case KindExperiment, KindManifest, KindEvent, KindDerived, KindPalette:
		return true
	default:
		return false
	}
}
