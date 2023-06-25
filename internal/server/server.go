package server

import (
	"net/http"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/server/handlers"
	_ "gocloud.dev/blob/fileblob"
)

type Server struct {
	Mux *http.ServeMux
}

func NewServer(ps photos.Photos, conf *config.Config) (*Server, error) {
	srv := &Server{
		Mux: http.NewServeMux(),
	}

	srv.Mux.Handle("/", handlers.IndexHandler{
		Photos: ps,
		Conf:   conf,
	})

	srv.Mux.Handle("/api/", handlers.APIHandler{
		Photos: ps,
		Conf:   conf,
	})

	return srv, nil
}
