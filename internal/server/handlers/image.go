package handlers

import (
	"net/http"

	"github.com/yklcs/panchro/storage"
)

type ImageHandler struct {
	Store storage.Storage
}

func (h ImageHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.Store.ServeHTTP(w, r)
}
