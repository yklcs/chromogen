package handlers

import (
	"net/http"

	"github.com/yklcs/chromogen/internal/config"
	"github.com/yklcs/chromogen/internal/photos"
	"github.com/yklcs/chromogen/internal/render"
)

type ChromogenHandler struct {
	Photos *photos.Photos
	Conf   *config.Config
}

func (h ChromogenHandler) Get(w http.ResponseWriter, r *http.Request) {
	render.RenderChromogen(w, h.Photos, h.Conf)
}
