package serve

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	serve "github.com/yklcs/panchro/serve/internal"
	"github.com/yklcs/panchro/storage"
)

type Server struct {
	// storepath string
	// dbpath    string
	conf   *config.Config
	photos *photos.Photos
	port   string
	router *chi.Mux
}

func NewServer(port, storepath, dbpath, confpath, s3url string) (*Server, error) {
	conf, err := config.ReadConfig(confpath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbpath)
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

	srv, _ := serve.NewRouter(ps, store, conf)

	return &Server{
		router: srv,
		conf:   conf,
		port:   port,
		photos: ps,
	}, nil
}

func (srv *Server) Serve() error {
	// err := srv.photos.Init()
	// if err != nil {
	// return err
	// }
	err := http.ListenAndServe(":"+srv.port, srv.router)
	return err
}

func (srv *Server) Close() error {
	return srv.photos.DB.Close()
}
