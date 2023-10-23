package handlers

import (
	"log"
	"net/http"

	"github.com/yklcs/chromogen/internal/config"
	"github.com/yklcs/chromogen/internal/photos"
	"github.com/yklcs/chromogen/internal/render"
)

type IndexHandler struct {
	Photos *photos.Photos
	Conf   *config.Config
}

func (h IndexHandler) Get(w http.ResponseWriter, r *http.Request) {
	err := render.RenderIndex(w, h.Photos, h.Conf)
	if err != nil {
		log.Println(err)
	}
}
