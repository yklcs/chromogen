package server

import (
	"net/http"

	"github.com/yklcs/panchro/internal/server/handlers"
	"github.com/yklcs/panchro/web"
)

func (s *Server) InitMux() error {
	s.Mux.Handle("/", http.FileServer(http.FS(web.Assets)))
	s.Mux.Handle("/images/", handlers.ImageHandler{
		Bucket: s.Bucket,
	})
	s.Mux.Handle("/images", handlers.ImagesHandler{
		Bucket: s.Bucket,
	})
	return nil
}
