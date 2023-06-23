package cmd

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/server"
	_ "gocloud.dev/blob/fileblob"
)

func Serve(args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	var outDir = flags.String("o", "dist", "output directory")
	var confPath = flags.String("c", "panchro.json", "configuration json file path")
	var concurrency = flags.Int("concurrency", 128, "configuration json file path")
	// var compressImages = flags.Bool("compress", true, "enable image compression")

	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: panchro serve [...flags] <input url>")
		fmt.Fprintln(flags.Output(), "Flags:")
		flags.PrintDefaults()
		fmt.Fprintln(flags.Output(), "Example: panchro serve  -o=output -c=config.json images")
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

	inUrl := flags.Args()[0]
	if !strings.HasPrefix(inUrl, "s3://") && !strings.HasPrefix(inUrl, "file://") {
		absInUrl, err := filepath.Abs(inUrl)
		if err != nil {
			return err
		}

		inUrlUrl := url.URL{
			Scheme: "file",
			Path:   absInUrl,
		}
		inUrl = inUrlUrl.String()
	}

	conf, err := config.ReadConfig(*confPath)
	if err != nil {
		return err
	}

	ps := photos.Photos{}

	err = ps.Read(inUrl, *outDir, *concurrency)
	if err != nil {
		return err
	}

	if ps.Len() == 0 {
		return errors.New("no images found in " + ps.BucketURL)
	}

	srv, _ := server.NewServer(ps, conf)
	http.ListenAndServe(":8000", srv.Mux)

	return nil
}
