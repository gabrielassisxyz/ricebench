module github.com/gabrielassisxyz/ricebench

// The patch version is pinned, not just the minor one. go1.26.2 ships four standard-library
// vulnerabilities that govulncheck traces straight through net/http from this binary's only
// entry point, so the gate is red on any toolchain older than this. Pinning here fixes it for
// every clone through toolchain switching, rather than depending on what each machine has
// installed. Raise it when a later patch is required; do not lower it to match a local Go.
go 1.26.6
