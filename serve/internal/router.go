package serve

import (
	"io/fs"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/yklcs/chromogen/internal/config"
	"github.com/yklcs/chromogen/internal/photos"
	"github.com/yklcs/chromogen/serve/internal/handlers"
	"github.com/yklcs/chromogen/storage"
)

func NewRouter(ps *photos.Photos, store storage.Storage, conf *config.Config) (*chi.Mux, error) {
	r := chi.NewRouter()

	themeFS := config.LoadTheme(conf)
	staticFS, _ := fs.Sub(themeFS, "static")
	staticHandler := handlers.StaticHandler{
		Handler: http.FileServer(http.FS(staticFS)),
		Conf:    conf,
	}

	theme, err := config.NewTheme(conf)
	if err != nil {
		return nil, err
	}
	photoHandler := handlers.PhotoHandler{Photos: ps, Conf: conf, Theme: theme}
	imageHandler := handlers.ImageHandler{Store: store}
	indexHandler := handlers.IndexHandler{
		Photos: ps,
		Conf:   conf,
		Theme:  theme,
	}

	r.Get(path.Join("/", conf.StaticDir, "*"), staticHandler.Get)
	r.Get("/{id}", photoHandler.Get)
	r.Get("/i/{id}.jpeg", imageHandler.Get)
	r.Get("/", indexHandler.Get)

	return r, nil
}
