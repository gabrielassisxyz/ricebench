package server

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"

	"github.com/gabrielassisxyz/ricebench/internal/fixture"
)

// galleryPath is the development-only visual review route. It is served by the Go
// process, never by the embedded frontend, so it cannot reach a participant build.
const galleryPath = "/dev/gallery"

// galleryHandler serves the development-only visual review surface: every fixture
// family rendered under a selected extreme palette, plus the required-role coverage
// view. It exists so a reviewer can see a visibility or hierarchy failure before real
// candidates do, and it never exposes internal identities to a participant flow.
func galleryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		writeGallery(w, fixture.PaletteID(r.URL.Query().Get("palette")))
	})
}

func writeGallery(w io.Writer, selected fixture.PaletteID) {
	palettes := fixture.ExtremePalettes()
	chosen := findExtremePalette(palettes, selected)

	writeHTML(w, "<!doctype html><html><head><meta charset=\"utf-8\"><title>RiceBench development gallery</title>")
	writeHTML(w, "<style>body{font-family:ui-monospace,monospace;margin:2rem;background:#111;color:#ddd}section{margin:1rem 0}.region,.block{border:1px solid #333;padding:.5rem;margin:.25rem 0}.run{white-space:pre-wrap}table{border-collapse:collapse}td,th{border:1px solid #333;padding:.25rem .5rem;text-align:left;vertical-align:top}</style>")
	writeHTML(w, "</head><body>")
	writeHTML(w, "<h1>Development gallery</h1>")
	writeHTML(w, "<p>Development-only inspection surface. It is not part of the participant flow and never appears in built participant assets.</p>")

	writeHTML(w, "<h2>Extreme palettes</h2><ul>")
	for _, palette := range palettes {
		mark := ""
		if chosen != nil && palette.ID == chosen.ID {
			mark = " (selected)"
		}
		writeHTML(w, "<li><a href=\"%s?palette=%s\">%s</a>%s: %s</li>",
			galleryPath, html.EscapeString(string(palette.ID)), html.EscapeString(string(palette.ID)), mark, html.EscapeString(palette.Description))
	}
	writeHTML(w, "</ul>")

	if chosen == nil {
		if selected != "" {
			writeHTML(w, "<p>Unknown palette %q. Pick one from the list above.</p>", html.EscapeString(string(selected)))
		}
		writeHTML(w, "</body></html>")
		return
	}

	colors, err := fixture.ResolvedColors(chosen.Palette)
	if err != nil {
		writeHTML(w, "<p>Palette %q does not resolve: %s</p>", html.EscapeString(string(chosen.ID)), html.EscapeString(err.Error()))
		writeHTML(w, "</body></html>")
		return
	}

	writeHTML(w, "<h2>%s</h2><p>%s</p>", html.EscapeString(string(chosen.ID)), html.EscapeString(chosen.Description))
	for _, set := range allGalleryFixtureSets() {
		writeGalleryFamily(w, set, colors)
	}

	writeGalleryCoverage(w, allGalleryFixtureSets()...)
	writeHTML(w, "</body></html>")
}

func writeGalleryFamily(w io.Writer, set fixture.FixtureSet, colors map[fixture.RoleID]string) {
	writeHTML(w, "<section class=\"family\" data-fixture-id=\"%s\"><h3>%s</h3>", html.EscapeString(string(set.ID)), html.EscapeString(string(set.ID)))
	for _, scene := range set.Scenes {
		writeHTML(w, "<section class=\"scene\" data-scene-id=\"%s\" data-scene-family=\"%s\">", html.EscapeString(string(scene.ID)), html.EscapeString(string(scene.Family)))
		for _, region := range scene.Regions {
			writeGalleryRegion(w, region, colors)
		}
		writeHTML(w, "</section>")
	}
	writeHTML(w, "</section>")
}

func writeGalleryRegion(w io.Writer, region fixture.Region, colors map[fixture.RoleID]string) {
	writeHTML(w, "<div class=\"region\" data-region-id=\"%s\" data-region-kind=\"%s\" data-semantic-state=\"%s\"%s%s%s style=\"%s\">",
		html.EscapeString(string(region.ID)),
		html.EscapeString(string(region.Kind)),
		html.EscapeString(string(region.State)),
		roleAttribute("data-role-background", region.Background),
		roleAttribute("data-role-foreground", region.Foreground),
		roleAttribute("data-role-border", region.Border),
		primitiveStyle(region.Background, region.Foreground, region.Border, colors),
	)
	for _, block := range region.Blocks {
		writeGalleryBlock(w, block, colors)
	}
	writeHTML(w, "</div>")
}

func writeGalleryBlock(w io.Writer, block fixture.ContentBlock, colors map[fixture.RoleID]string) {
	writeHTML(w, "<div class=\"block\" data-block-id=\"%s\" data-content-kind=\"%s\" data-semantic-state=\"%s\"%s%s%s style=\"%s\">",
		html.EscapeString(string(block.ID)),
		html.EscapeString(string(block.Kind)),
		html.EscapeString(string(block.State)),
		roleAttribute("data-role-background", block.Background),
		roleAttribute("data-role-foreground", block.Foreground),
		roleAttribute("data-role-border", block.Border),
		primitiveStyle(block.Background, block.Foreground, block.Border, colors),
	)
	if len(block.Runs) > 0 {
		for _, run := range block.Runs {
			writeGalleryRun(w, run, colors)
		}
	} else {
		writeHTML(w, "%s", html.EscapeString(block.Text))
	}
	writeHTML(w, "</div>")
}

func writeGalleryRun(w io.Writer, run fixture.ContentRun, colors map[fixture.RoleID]string) {
	writeHTML(w, "<span class=\"run\" data-semantic-state=\"%s\"%s%s style=\"%s\">%s</span>",
		html.EscapeString(string(run.State)),
		roleAttribute("data-role-background", run.Background),
		roleAttribute("data-role-foreground", run.Foreground),
		primitiveStyle(run.Background, run.Foreground, "", colors),
		html.EscapeString(run.Text),
	)
}

func writeGalleryCoverage(w io.Writer, sets ...fixture.FixtureSet) {
	writeHTML(w, "<h2>Coverage</h2><table><tr><th>role</th><th>references</th></tr>")
	for _, entry := range fixture.CoverageView(sets...) {
		references := make([]string, 0, len(entry.References))
		for _, reference := range entry.References {
			references = append(references, fmt.Sprintf("%s (%s)", reference.Path, reference.Level))
		}
		writeHTML(w, "<tr><td>%s</td><td>%s</td></tr>",
			html.EscapeString(string(entry.Role)),
			html.EscapeString(strings.Join(references, ", ")),
		)
	}
	writeHTML(w, "</table>")
}

func findExtremePalette(palettes []fixture.ExtremePalette, id fixture.PaletteID) *fixture.ExtremePalette {
	for index := range palettes {
		if palettes[index].ID == id {
			return &palettes[index]
		}
	}
	return nil
}

func allGalleryFixtureSets() []fixture.FixtureSet {
	return []fixture.FixtureSet{
		fixture.TerminalAgentFixtureSet(),
		fixture.CodeDiffFixtureSet(),
		fixture.DesktopShellFixtureSet(),
		fixture.ReadingMonitoringFixture(),
	}
}

func roleAttribute(name string, role fixture.RoleID) string {
	if role == "" {
		return ""
	}
	return fmt.Sprintf(" %s=\"%s\"", name, html.EscapeString(string(role)))
}

func primitiveStyle(background, foreground, border fixture.RoleID, colors map[fixture.RoleID]string) string {
	var parts []string
	if background != "" {
		parts = append(parts, "background-color:"+colors[background])
	}
	if foreground != "" {
		parts = append(parts, "color:"+colors[foreground])
	}
	if border != "" {
		parts = append(parts, "border-color:"+colors[border])
	}
	return strings.Join(parts, ";")
}

// writeHTML writes to the gallery output. A write error here means the client went
// away mid-response and there is nothing to recover, so the error is deliberately
// ignored rather than threaded through every render helper.
func writeHTML(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
