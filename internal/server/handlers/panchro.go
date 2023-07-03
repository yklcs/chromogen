package handlers

import (
	"net/http"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/render"
)

type PanchroHandler struct {
	Photos *photos.Photos
	Conf   *config.Config
}

func (h PanchroHandler) Get(w http.ResponseWriter, r *http.Request) {
	render.RenderPanchro(w, h.Photos, h.Conf)
}
