package panchro

import (
	"net/http"
	"strings"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/server"
	"github.com/yklcs/panchro/storage"
	bolt "go.etcd.io/bbolt"
)

type Server struct {
	// storepath string
	// dbpath    string
	conf   *config.Config
	photos *photos.Photos
	port   string
	server *server.Server
}

func NewServer(port, storepath, dbpath, confpath, s3url string) (*Server, error) {
	conf, err := config.ReadConfig(confpath)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		return nil, err
	}

	ps := &photos.Photos{DB: db}
	ps.Init()

	var store storage.Storage
	if strings.HasPrefix(storepath, "s3://") {
		s3path, _ := strings.CutPrefix(storepath, "s3://")
		bucket, prefix, _ := strings.Cut(s3path, "/")
		store, err = storage.NewS3Storage(bucket, prefix, s3url)
	} else {
		store, err = storage.NewLocalStorage(storepath)
	}
	if err != nil {
		return nil, err
	}

	srv, _ := server.NewServer(ps, store, conf)

	return &Server{
		server: srv,
		conf:   conf,
		port:   port,
		photos: ps,
	}, nil

	return nil, nil
}

func (srv *Server) Serve() error {
	err := srv.photos.Init()
	if err != nil {
		return err
	}
	err = http.ListenAndServe(":"+srv.port, srv.server.Router)
	return err
}

func (srv *Server) Close() error {
	return srv.photos.DB.Close()
}
