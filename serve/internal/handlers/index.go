package handlers

import (
	"log"
	"net/http"

	"github.com/yklcs/chromogen/internal/config"
	"github.com/yklcs/chromogen/internal/photos"
)

type IndexHandler struct {
	Photos *photos.Photos
	Conf   *config.Config
	Theme  *config.Theme
}

func (h IndexHandler) Get(w http.ResponseWriter, r *http.Request) {
	err := h.Theme.Render(w, "index",
		config.ThemeData{Photos: h.Photos, Config: h.Conf})
	if err != nil {
		log.Println(err)
	}
}
