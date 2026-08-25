package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestGalleryListsPalettesWhenNoneSelected(t *testing.T) {
	handler, err := New(builtAssets())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	body := getGallery(t, handler, galleryPath)
	for _, id := range []string{"extreme-flat-hierarchy", "extreme-collapsed-surfaces", "extreme-invisible-cursor-selection", "extreme-semantic-collision", "extreme-ansi-bright-collapse"} {
		if !strings.Contains(body, id) {
			t.Fatalf("gallery did not list palette %q", id)
		}
	}
}

func TestGalleryRendersFamilyUnderSelectedPalette(t *testing.T) {
	handler, err := New(builtAssets())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	body := getGallery(t, handler, galleryPath+"?palette=extreme-flat-hierarchy")

	// The response must carry the fixture structure, not just a 200.
	for _, marker := range []string{
		`data-scene-id="terminal-shell"`,
		`data-region-id="shell-frame"`,
		`data-role-background="terminal-background"`,
		`data-block-id="shell-prompt"`,
	} {
		if !strings.Contains(body, marker) {
			t.Fatalf("gallery response does not carry %q", marker)
		}
	}

	// The response must carry the fixture rendered under the requested palette, which
	// means asserting on what distinguishes that palette rather than on any resolved
	// color. Every extreme palette shares a base where background is #1a1a1a, so a
	// response serving the wrong palette still carries it. Flat hierarchy is the one
	// that collapses surface and surface-elevated onto the background, so the base
	// values for those two roles must be absent.
	if !strings.Contains(body, "background-color:#1a1a1a") {
		t.Fatal("gallery response does not carry the resolved palette color")
	}
	for _, collapsed := range []string{"#242424", "#2e2e2e"} {
		if strings.Contains(body, collapsed) {
			t.Fatalf("gallery served a palette that keeps %s: flat hierarchy collapses it", collapsed)
		}
	}
}

func TestGalleryRejectsUnknownPalette(t *testing.T) {
	handler, err := New(builtAssets())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	body := getGallery(t, handler, galleryPath+"?palette=does-not-exist")
	if !strings.Contains(body, "Unknown palette") {
		t.Fatal("gallery did not report the unknown palette")
	}
}

func TestGalleryStructureIsPaletteIndependent(t *testing.T) {
	handler, err := New(builtAssets())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	flat := getGallery(t, handler, galleryPath+"?palette=extreme-flat-hierarchy")
	collapsed := getGallery(t, handler, galleryPath+"?palette=extreme-collapsed-surfaces")

	if flat == collapsed {
		t.Fatal("palette switch produced identical output; colors did not change")
	}
	if galleryStructure(flat) != galleryStructure(collapsed) {
		t.Fatal("palette switch changed structure or role-reference identifiers")
	}
}

func TestGalleryAbsentFromParticipantAssets(t *testing.T) {
	handler, err := New(builtAssets())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	index := getGallery(t, handler, "/")
	if strings.Contains(index, "Development gallery") || strings.Contains(index, galleryPath) {
		t.Fatal("participant assets expose the development gallery")
	}

	gallery := getGallery(t, handler, galleryPath)
	if !strings.Contains(gallery, "Development gallery") {
		t.Fatal("gallery route did not serve the development gallery")
	}
}

func getGallery(t *testing.T, handler http.Handler, path string) string {
	t.Helper()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("GET %s: status %d", path, recorder.Code)
	}
	body, _ := io.ReadAll(recorder.Body)
	return string(body)
}

var galleryStylePattern = regexp.MustCompile(` style="[^"]*"`)

// galleryStructure returns the family and coverage sections with resolved colors
// removed, so two renders can be compared for structure alone.
func galleryStructure(body string) string {
	start := strings.Index(body, `<section class="family"`)
	if start < 0 {
		return body
	}
	return galleryStylePattern.ReplaceAllString(body[start:], "")
}
