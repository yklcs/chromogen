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
	bolt "go.etcd.io/bbolt"
)

func Build(args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	out := flags.String("o", "dist", "output directory")
	confPath := flags.String("c", "panchro.json", "configuration json file path")

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

	dbpath := path.Join(*out, "tmp.db")
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	defer os.RemoveAll(dbpath)

	ps := &photos.Photos{DB: db}
	ps.Init()
	err = ps.ProcessFS(in, *out, true, 2048, 75)
	if err != nil {
		return err
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

	themeFS := config.LoadTheme(conf)
	staticFS, _ := fs.Sub(themeFS, "static")
	staticDir := path.Join(*out, conf.StaticDir)
	err = os.MkdirAll(staticDir, 0755)
	if err != nil {
		return err
	}
	err = render.CopyFS(staticFS, staticDir)
	if err != nil {
		return err
	}

	for _, id := range ps.IDs() {
		dir := path.Join(*out, id)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		imageHTML, err := os.Create(path.Join(dir, "index.html"))
		if err != nil {
			return err
		}
		defer imageHTML.Close()

		p, _ := ps.Get(id)
		err = render.RenderPhoto(imageHTML, &p, conf)
		if err != nil {
			return err
		}
	}

	return nil
}
