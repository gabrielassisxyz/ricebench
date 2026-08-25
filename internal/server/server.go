// Package server wires the HTTP surface RiceBench serves on a local listener.
package server

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"
)

// ErrAssetsMissing reports a binary built without its frontend. It is returned rather than
// served as an empty tree because a server that answers 404 for every path looks like a
// routing bug and sends the reader looking in the wrong place.
var ErrAssetsMissing = errors.New("frontend assets missing: run bin/build-web, then rebuild")

// New returns the handler serving the embedded frontend.
func New(assets fs.FS) (http.Handler, error) {
	if _, err := fs.Stat(assets, "index.html"); err != nil {
		return nil, ErrAssetsMissing
	}

	mux := http.NewServeMux()
	mux.Handle(galleryPath, galleryHandler())
	mux.Handle("/", http.FileServerFS(assets))
	return mux, nil
}

// StartupNotice describes where the server is reachable.
//
// A non-loopback listener gets an explicit warning because RiceBench has no authentication:
// any client that can reach the port can read the experiment and submit judgments. The
// exposure is an opt-in the operator has to be able to see, and nothing else in the process
// will tell them.
func StartupNotice(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Sprintf("RiceBench listening on %s\n", addr)
	}

	var notice strings.Builder
	fmt.Fprintf(&notice, "RiceBench listening on http://%s/\n", net.JoinHostPort(reachableHost(host), port))
	if !isLoopback(host) {
		notice.WriteString(
			"WARNING: this listener is not loopback and the session is unauthenticated.\n" +
				"Every client that can reach this port can read the experiment and submit judgments.\n" +
				"RiceBench changes no firewall rule; network reachability stays a machine concern.\n")
	}
	return notice.String()
}

// reachableHost turns a wildcard bind into something a browser can open. An empty host or
// an unspecified address means every interface, and "http://:7391/" is not a usable URL.
func reachableHost(host string) string {
	if host == "" {
		return "localhost"
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsUnspecified() {
		return "localhost"
	}
	return host
}

func isLoopback(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
