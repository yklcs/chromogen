package serve

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Cmd(args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	s3url := flags.String("s3url", "", "S3 URL root, use if S3 is behind CDN")
	confpath := flags.String("c", "chromogen.json", "configuration json file path")
	port := flags.String("p", "8000", "port")

	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: chromogen serve [...flags] <input url>")
		fmt.Fprintln(flags.Output(), "Flags:")
		flags.PrintDefaults()
		fmt.Fprintln(flags.Output(), "Example: chromogen serve  -o=output -c=config.json images")
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

	inpath := flags.Args()[0]
	srv, err := NewServer(*port, inpath, "chromogen", *confpath, *s3url)
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		srv.Close()
		os.Exit(1)
	}()

	err = srv.Serve()
	if err != nil {
		return err
	}

	return nil
}
