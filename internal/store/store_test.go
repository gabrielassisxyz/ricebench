package store

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveReturnsExplicitRoot(t *testing.T) {
	dir := t.TempDir()

	root, err := Resolve(dir)
	if err != nil {
		t.Fatalf("Resolve(%q): %v", dir, err)
	}
	want, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("filepath.Abs(%q): %v", dir, err)
	}
	if root.Dir() != filepath.Clean(want) {
		t.Fatalf("Dir() = %q, want %q", root.Dir(), filepath.Clean(want))
	}
}

func TestResolveRejectsOmittedRoot(t *testing.T) {
	_, err := Resolve("")
	if !errors.Is(err, ErrRootOmitted) {
		t.Fatalf("Resolve(\"\"): got %v, want ErrRootOmitted", err)
	}
}

func TestResolveRejectsMissingRoot(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "missing")

	_, err := Resolve(dir)
	if !errors.Is(err, ErrRootMissing) {
		t.Fatalf("Resolve(%q): got %v, want ErrRootMissing", dir, err)
	}
}

func TestResolveRejectsInaccessibleRoot(t *testing.T) {
	path := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := Resolve(path)
	if !errors.Is(err, ErrRootInaccessible) {
		t.Fatalf("Resolve(%q): got %v, want ErrRootInaccessible", path, err)
	}
}

func TestPathResolvesEveryKindInsideRoot(t *testing.T) {
	root, err := Resolve(t.TempDir())
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	kinds := []struct {
		kind Kind
		name string
	}{
		{KindExperiment, "experiment.json"},
		{KindManifest, "manifest.json"},
		{KindEvent, "events/0001.jsonl"},
		{KindDerived, "derived/profile.json"},
		{KindPalette, "palette.json"},
	}
	for _, tc := range kinds {
		got, err := root.Path(tc.kind, tc.name)
		if err != nil {
			t.Fatalf("Path(%s, %q): %v", tc.kind, tc.name, err)
		}
		want := filepath.Join(root.Dir(), filepath.FromSlash(tc.name))
		if got != want {
			t.Fatalf("Path(%s, %q) = %q, want %q", tc.kind, tc.name, got, want)
		}
		if !strings.HasPrefix(got, root.Dir()+string(filepath.Separator)) {
			t.Fatalf("Path(%s, %q) = %q does not lie under root %q", tc.kind, tc.name, got, root.Dir())
		}
	}
}

func TestPathRejectsTraversal(t *testing.T) {
	root, err := Resolve(t.TempDir())
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	for _, name := range []string{"..", "../outside", "a/../../outside"} {
		got, err := root.Path(KindExperiment, name)
		if !errors.Is(err, ErrPathEscapes) {
			t.Fatalf("Path(KindExperiment, %q) = %q, %v; want ErrPathEscapes", name, got, err)
		}
	}
}

func TestPathRejectsInvalidNames(t *testing.T) {
	root, err := Resolve(t.TempDir())
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	cases := []struct {
		kind Kind
		name string
	}{
		{KindManifest, ""},
		{KindManifest, "."},
		{KindManifest, "a\x00b"},
		{"unknown", "manifest.json"},
	}
	for _, tc := range cases {
		got, err := root.Path(tc.kind, tc.name)
		if !errors.Is(err, ErrPathInvalid) {
			t.Fatalf("Path(%q, %q) = %q, %v; want ErrPathInvalid", tc.kind, tc.name, got, err)
		}
	}
}

func TestPathRejectsAbsoluteName(t *testing.T) {
	root, err := Resolve(t.TempDir())
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	got, err := root.Path(KindEvent, filepath.Join(root.Dir(), "events.log"))
	if !errors.Is(err, ErrPathInvalid) {
		t.Fatalf("Path(KindEvent, absolute) = %q, %v; want ErrPathInvalid", got, err)
	}
}
