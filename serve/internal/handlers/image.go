package handlers

import (
	"net/http"

	"github.com/yklcs/chromogen/storage"
)

type ImageHandler struct {
	Store storage.Storage
}

func (h ImageHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.Store.ServeHTTP(w, r)
}
