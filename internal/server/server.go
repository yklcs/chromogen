package server

import (
	"io/fs"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/server/handlers"
	"github.com/yklcs/panchro/storage"
	"github.com/yklcs/panchro/web"
	_ "gocloud.dev/blob/fileblob"
)

type Server struct {
	Router *chi.Mux
}

func NewServer(ps *photos.Photos, store storage.Storage, conf *config.Config) (*Server, error) {
	r := chi.NewRouter()

	photosHandler := handlers.PhotosHandler{
		Photos: ps,
		Store:  store,
	}
	staticFs, _ := fs.Sub(web.Content, "static")
	staticHandler := handlers.StaticHandler{
		Handler: http.FileServer(http.FS(staticFs)),
		Conf:    conf,
	}
	photoHandler := handlers.PhotoHandler{Photos: ps, Conf: conf}
	imageHandler := handlers.ImageHandler{Store: store}
	indexHandler := handlers.IndexHandler{
		Photos: ps,
		Conf:   conf,
	}
	panchroHandler := handlers.PanchroHandler{
		Photos: ps,
		Conf:   conf,
	}

	auth := Auth("secret")
	r.With(auth).Post("/photos", photosHandler.Post)
	r.With(auth).Delete("/photos/{id}", photosHandler.Delete)
	r.Get("/photos", photosHandler.Get)
	r.Get(path.Join("/", conf.StaticDir, "*"), staticHandler.Get)
	r.Get("/{id}", photoHandler.Get)
	r.Get("/{id}.jpg", imageHandler.Get)
	r.Get("/", indexHandler.Get)
	r.Get("/panchro", panchroHandler.Get)

	return &Server{
		Router: r,
	}, nil
}
