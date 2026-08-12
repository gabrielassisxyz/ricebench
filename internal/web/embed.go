// Package web carries the built frontend assets that the server hands to the browser.
package web

import (
	"embed"
	"io/fs"
)

// The `all:` prefix is load-bearing. Plain `dist` ignores names starting with a dot, so on
// a clone whose frontend was never built the only file present is dist/.gitkeep, go:embed
// finds no match, and every Go command fails before anyone can install Node.
//
//go:embed all:dist
var assets embed.FS

// Dist returns the built frontend rooted at the directory bin/build-web writes.
func Dist() (fs.FS, error) {
	return fs.Sub(assets, "dist")
}
