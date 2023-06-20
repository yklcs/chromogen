package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yklcs/panchro/cmd"
)

func usage() {
	fmt.Print(`panchro [command]
Usage:
	panchro build		build panchro site
	panchro serve		serve panchro site

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
		run = cmd.Build
	case "serve":
		run = cmd.Serve
	case "help":
		usage()
		os.Exit(0)
	default:
		usage()
		os.Exit(1)
	}

	if err := run(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
