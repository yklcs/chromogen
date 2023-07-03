package handlers

import (
	"net/http"

	"github.com/yklcs/panchro/storage"
)

type ImageHandler struct {
	Store storage.Storage
}

func (h ImageHandler) Get(w http.ResponseWriter, r *http.Request) {
	// id := chi.URLParam(r, "id")
	h.Store.ServeHTTP(w, r)
}
