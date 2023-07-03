package cmd

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/dgraph-io/badger/v3"
	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/server"
	"github.com/yklcs/panchro/storage"
	_ "gocloud.dev/blob/fileblob"
)

func Serve(args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	storename := flags.String("s", "panchro", "storage (local/s3)")
	confPath := flags.String("c", "panchro.json", "configuration json file path")
	var port = flags.String("p", "8000", "port")

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

	if len(flags.Args()) != 0 {
		flags.Usage()
		return errors.New("wrong number of arguments")
	}

	conf, err := config.ReadConfig(*confPath)
	if err != nil {
		return err
	}

	db, err := badger.Open(badger.DefaultOptions(path.Join(*storename, "panchro.db")))
	if err != nil {
		return err
	}

	store, err := storage.NewLocalStorage(*storename)
	if err != nil {
		return err
	}
	ps := photos.NewPhotos(db)
	srv, _ := server.NewServer(&ps, store, conf)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		db.Close()
		os.Exit(1)
	}()

	http.ListenAndServe(":"+*port, srv.Router)

	return nil
}
