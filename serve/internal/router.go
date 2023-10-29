package serve

import (
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yklcs/chromogen/internal/config"
	"github.com/yklcs/chromogen/internal/photos"
	"github.com/yklcs/chromogen/serve/internal/handlers"
	"github.com/yklcs/chromogen/storage"
)

func NewRouter(ps *photos.Photos, store storage.Storage, conf *config.Config) (*chi.Mux, error) {
	r := chi.NewRouter()

	theme, err := config.NewTheme(conf)
	if err != nil {
		return nil, err
	}

	staticHandler := handlers.StaticHandler{
		Handler: http.FileServer(http.FS(theme.StaticFS)),
		Conf:    conf,
	}

	photoHandler := handlers.PhotoHandler{Photos: ps, Conf: conf, Theme: theme}
	imageHandler := handlers.ImageHandler{Store: store}
	indexHandler := handlers.IndexHandler{
		Photos: ps,
		Conf:   conf,
		Theme:  theme,
	}

	r.Use(middleware.StripSlashes)
	r.Get(path.Join("/", conf.StaticDir, "*"), staticHandler.Get)
	r.Get("/{id}", photoHandler.Get)
	r.Get("/i/*", imageHandler.Get)
	r.Get("/", indexHandler.Get)

	return r, nil
}
