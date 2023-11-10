package handlers

import (
	"log"
	"net/http"

	"github.com/yklcs/chromogen/internal/config"
	"github.com/yklcs/chromogen/internal/photos"
	"github.com/yklcs/chromogen/internal/theme"
)

type IndexHandler struct {
	Photos *photos.Photos
	Conf   *config.Config
	Theme  *theme.Theme
}

func (h IndexHandler) Get(w http.ResponseWriter, r *http.Request) {
	err := h.Theme.Render(w, "index",
		theme.ThemeData{Photos: h.Photos, Config: h.Conf})
	if err != nil {
		log.Println(err)
	}
}
