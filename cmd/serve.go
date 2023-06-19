package cmd

import (
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/server"
	_ "gocloud.dev/blob/fileblob"
)

func Serve(args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	var inUrl = flags.String("i", ".", "input")
	var outDir = flags.String("o", "dist", "output directory")
	var confPath = flags.String("c", "panchro.json", "configuration json file path")
	// var compressImages = flags.Bool("compress", true, "enable image compression")

	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: panchro build")
		fmt.Fprintln(flags.Output(), "Flags:")
		flags.PrintDefaults()
	}

	err := flags.Parse(args)
	if err != nil {
		return err
	}

	if len(flags.Args()) > 0 {
		flags.Usage()
		return errors.New("unknown arguments")
	}

	conf, err := config.ReadConfig(*confPath)
	if err != nil {
		return err
	}

	ps := photos.Photos{}

	err = ps.Read(*inUrl, *outDir)
	if err != nil {
		return err
	}

	// Reverse imgs
	// for i, j := 0, len(imgs)-1; i < j; i, j = i+1, j-1 {
	// 	imgs[i], imgs[j] = imgs[j], imgs[i]
	// }

	if ps.Len() == 0 {
		return errors.New("no images found in " + *inUrl)
	}

	srv, _ := server.NewServer(ps, conf)
	http.ListenAndServe(":8000", srv.Mux)

	return nil
}
