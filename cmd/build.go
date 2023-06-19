package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/render"
	"github.com/yklcs/panchro/web"
	_ "gocloud.dev/blob/fileblob"
)

func Build(args []string) error {
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

	ps := photos.NewPhotos(*outDir)

	err = os.RemoveAll(ps.Dir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(ps.Dir, 0755)
	if err != nil {
		return err
	}

	err = ps.Read(*inUrl, ps.Dir)
	if err != nil {
		return err
	}

	if ps.Len() == 0 {
		return errors.New("no images found in " + *inUrl)
	}

	indexHTML, err := os.Create(path.Join(*outDir, "index.html"))
	if err != nil {
		return err
	}

	err = render.RenderIndex(indexHTML, ps, conf)
	if err != nil {
		return err
	}

	static, _ := fs.Sub(web.Content, "static")
	err = render.CopyFS(static, *outDir)
	if err != nil {
		return err
	}

	for i := 0; i < ps.Len(); i++ {
		dir := path.Join(*outDir, ps.Get(i).ID)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		imageHTML, err := os.Create(path.Join(dir, "index.html"))
		if err != nil {
			return err
		}

		err = render.RenderPhoto(imageHTML, *ps.Get(i), conf)
		if err != nil {
			return err
		}
	}

	return nil
}
