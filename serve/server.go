package serve

import (
	"database/sql"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/yklcs/chromogen/internal/config"
	"github.com/yklcs/chromogen/internal/photos"
	serve "github.com/yklcs/chromogen/serve/internal"
	"github.com/yklcs/chromogen/storage"
)

type Server struct {
	conf   *config.Config
	photos *photos.Photos
	port   string
	router *chi.Mux
}

func NewServer(port, inpath, storepath, confpath, s3url string) (*Server, error) {
	conf, err := config.ReadConfig(confpath)
	if err != nil {
		return nil, err
	}

	store, err := storage.NewLocalStorage(storepath, "i")
	if err != nil {
		return nil, err
	}

	dbpath := path.Join(storepath, "chromogen.db")
	db, err := sql.Open("sqlite", dbpath)
	if err != nil {
		return nil, err
	}

	ps := &photos.Photos{DB: db}
	ps.Init()
	ps.LoadFiles([]string{inpath}, store)

	srv, _ := serve.NewRouter(ps, store, conf)

	return &Server{
		router: srv,
		conf:   conf,
		port:   port,
		photos: ps,
	}, nil
}

func (srv *Server) Serve() error {
	err := http.ListenAndServe(":"+srv.port, srv.router)
	return err
}

func (srv *Server) Close() error {
	return srv.photos.DB.Close()
}
