package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yklcs/chromogen/build"
	"github.com/yklcs/chromogen/serve"
)

func usage() {
	fmt.Print(`Usage: chromogen <command>
	chromogen build		build chromogen site
	chromogen serve		serve chromogen site (not for production)
	`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var run func([]string) error
	switch strings.ToLower(os.Args[1]) {
	case "build":
		run = build.Cmd
	case "serve":
		run = serve.Cmd
	case "help":
		usage()
		os.Exit(0)
	default:
		usage()
		os.Exit(1)
	}

	if err := run(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		fmt.Fprintf(os.Stderr, `run "chromogen %s -h" for help
`, os.Args[1])
		os.Exit(1)
	}
}
