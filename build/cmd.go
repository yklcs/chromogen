package build

import (
	"errors"
	"flag"
	"fmt"
)

func Cmd(args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	out := flags.String("o", "dist", "output directory")
	confpath := flags.String("c", "chromogen.json", "configuration json file path")

	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: chromogen build [...flags] <input url>")
		fmt.Fprintln(flags.Output(), "Flags:")
		flags.PrintDefaults()
		fmt.Fprintln(flags.Output(), "Example: chromogen build -o=output -c=config.json images")
		fmt.Fprintln(flags.Output())
	}

	err := flags.Parse(args)
	if err != nil {
		return err
	}

	if len(flags.Args()) == 0 {
		flags.Usage()
		return errors.New("wrong number of arguments")
	}

	in := flags.Args()

	ssg, err := NewStaticSiteGenerator(*out, *confpath)
	if err != nil {
		return err
	}
	err = ssg.Build(in)

	return err
}
