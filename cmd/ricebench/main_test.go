package main

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/gabrielassisxyz/ricebench/internal/store"
)

func TestRunRequiresDataDir(t *testing.T) {
	err := run("127.0.0.1:7391", "")
	if !errors.Is(err, store.ErrRootOmitted) {
		t.Fatalf("run with omitted data dir: got %v, want ErrRootOmitted", err)
	}
}

func TestRunRejectsMissingDataDir(t *testing.T) {
	err := run("127.0.0.1:7391", filepath.Join(t.TempDir(), "missing"))
	if !errors.Is(err, store.ErrRootMissing) {
		t.Fatalf("run with missing data dir: got %v, want ErrRootMissing", err)
	}
}
