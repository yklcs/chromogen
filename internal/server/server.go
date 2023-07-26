package server

import (
	"crypto/rand"
	"fmt"
	"io/fs"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/server/handlers"
	"github.com/yklcs/panchro/internal/utils"
	"github.com/yklcs/panchro/storage"
)

type Server struct {
	Router *chi.Mux
}

func NewServer(ps *photos.Photos, store storage.Storage, conf *config.Config) (*Server, error) {
	r := chi.NewRouter()

	randbytes := make([]byte, 32)
	rand.Read(randbytes)
	token := utils.Base58Encode(randbytes)
	fmt.Println(token)

	photosHandler := handlers.PhotosHandler{
		Photos: ps,
		Store:  store,
	}

	theme := config.LoadTheme(conf)
	staticFS, _ := fs.Sub(theme, "static")
	staticHandler := handlers.StaticHandler{
		Handler: http.FileServer(http.FS(staticFS)),
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

	auth := Auth(token)
	r.With(auth).Post("/photos", photosHandler.Post)
	r.With(auth).Delete("/photos/{id}", photosHandler.Delete)
	r.Get("/photos", photosHandler.GetAll)
	r.Get("/photos/{id}", photosHandler.Get)
	r.Get(path.Join("/", conf.StaticDir, "*"), staticHandler.Get)
	r.Get("/{id}", photoHandler.Get)
	r.Get("/{id}.jpg", imageHandler.Get)
	r.Get("/", indexHandler.Get)

	r.Route("/panchro", func(r chi.Router) {
		r.Use(AuthPage(token, conf))
		r.Get("/*", panchroHandler.Get)
	})

	return &Server{
		Router: r,
	}, nil
}
