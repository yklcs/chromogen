package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/render"
)

type PhotoHandler struct {
	Photos *photos.Photos
	Conf   *config.Config
}

func (h PhotoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.Photos.Get(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	render.RenderPhoto(w, &p, h.Conf)
}
