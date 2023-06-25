package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
)

type APIHandler struct {
	Photos photos.Photos
	Conf   *config.Config
}

func (h APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(h.Photos)
	}
	w.Header().Set("Content-Type", "application/json")
}
