package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yklcs/chromogen/internal/config"
	"github.com/yklcs/chromogen/internal/photos"
	"github.com/yklcs/chromogen/internal/theme"
)

type PhotoHandler struct {
	Photos *photos.Photos
	Conf   *config.Config
	Theme  *theme.Theme
}

func (h PhotoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.Photos.Get(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err = h.Theme.Render(w, "photo",
		theme.ThemeData{Photo: p, Config: h.Conf})
	if err != nil {
		log.Println(err)
	}
}
