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
	out := flags.String("o", "dist", "output directory")
	confPath := flags.String("c", "panchro.json", "configuration json file path")
	// var concurrency = flags.Int("concurrency", 128, "concurrency limit, edit depending on memory usage")
	compress := flags.Bool("compress", true, "enable image compression")

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

	conf, err := config.ReadConfig(*confPath)
	if err != nil {
		return errors.New(err.Error() + "\nspecify config file location with the -c flag")
	}

	err = os.MkdirAll(*out, 0755)
	if err != nil {
		return err
	}

	ps := photos.NewPhotos()
	err = ps.ProcessFS(in, *out, 2048, 75)
	if err != nil {
		return err
	}

	if ps.Len() == 0 {
		return errors.New("no images found in " + in)
	}

	if *compress {
		// err = ps.Compress()
		// if err != nil {
		// return err
		// }
	}

	indexHTML, err := os.Create(path.Join(*out, "index.html"))
	if err != nil {
		return err
	}
	defer indexHTML.Close()

	err = render.RenderIndex(indexHTML, ps, conf)
	if err != nil {
		return err
	}

	static, _ := fs.Sub(web.Content, "static")
	err = render.CopyFS(static, *out)
	if err != nil {
		return err
	}

	for i := 0; i < ps.Len(); i++ {
		dir := path.Join(*out, ps.Get(i).ID)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		imageHTML, err := os.Create(path.Join(dir, "index.html"))
		if err != nil {
			return err
		}

		defer imageHTML.Close()

		err = render.RenderPhoto(imageHTML, *ps.Get(i), conf)
		if err != nil {
			return err
		}
	}

	return nil
}
