// Command ricebench runs the local preference elicitation server and serves its frontend.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gabrielassisxyz/ricebench/internal/server"
	"github.com/gabrielassisxyz/ricebench/internal/web"
)

// version is stamped by the release build; a source build reports the placeholder.
var version = "dev"

func main() {
	addr := flag.String("addr", "127.0.0.1:7391",
		"address to listen on; a non-loopback value exposes the unauthenticated session to the network")
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	if err := run(*addr); err != nil {
		log.Fatalf("ricebench: %v", err)
	}
}

func run(addr string) error {
	assets, err := web.Dist()
	if err != nil {
		return err
	}

	handler, err := server.New(assets)
	if err != nil {
		return err
	}

	fmt.Print(server.StartupNotice(addr))
	return http.ListenAndServe(addr, handler)
}
