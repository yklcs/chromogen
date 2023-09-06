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
	storepath := flags.String("s", "panchro", "photo storage path, use s3://... for S3")
	s3url := flags.String("s3url", "", "S3 URL root, use if S3 is behind CDN")
	dbpath := flags.String("db", "panchro.db", "db path")
	confpath := flags.String("c", "panchro.json", "configuration json file path")
	port := flags.String("p", "8000", "port")

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

	srv, err := NewServer(*port, *storepath, *dbpath, *confpath, *s3url)
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

	// conf, err := config.ReadConfig(*confpath)
	// if err != nil {
	// 	return err
	// }

	// db, err := bolt.Open(*dbpath, 0600, nil)
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()

	// var store storage.Storage
	// if strings.HasPrefix(*storepath, "s3://") {
	// 	s3path, _ := strings.CutPrefix(*storepath, "s3://")
	// 	bucket, prefix, _ := strings.Cut(s3path, "/")
	// 	store, err = storage.NewS3Storage(bucket, prefix, *s3url)
	// } else {
	// 	store, err = storage.NewLocalStorage(*storepath)
	// }
	// if err != nil {
	// 	return err
	// }
	// ps := photos.Photos{DB: db}
	// ps.Init()
	// srv, _ := server.NewServer(&ps, store, conf)

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-c
	// 	db.Close()
	// 	os.Exit(1)
	// }()

	// http.ListenAndServe(":"+*port, srv.Router)

	return nil
}
