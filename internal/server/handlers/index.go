package handlers

import (
	"net/http"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/render"
)

type IndexHandler struct {
	Photos *photos.Photos
	Conf   *config.Config
}

func (h IndexHandler) Get(w http.ResponseWriter, r *http.Request) {
	render.RenderIndex(w, h.Photos, h.Conf)
}
