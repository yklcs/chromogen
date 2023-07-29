package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/yklcs/panchro"
)

func build(args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	out := flags.String("o", "dist", "output directory")
	confpath := flags.String("c", "panchro.json", "configuration json file path")

	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: panchro build [...flags] <input url>")
		fmt.Fprintln(flags.Output(), "Flags:")
		flags.PrintDefaults()
		fmt.Fprintln(flags.Output(), "Example: panchro build -o=output -c=config.json images")
		fmt.Fprintln(flags.Output())
	}

	err := flags.Parse(args)
	if err != nil {
		return err
	}

	if len(flags.Args()) != 1 {
		flags.Usage()
		return errors.New("wrong number of arguments")
	}

	in := flags.Args()[0]

	ssg, err := panchro.NewStaticSiteGenerator(in, *out, *confpath)
	if err != nil {
		return err
	}
	err = ssg.Build()

	return err
}
