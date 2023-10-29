package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yklcs/chromogen/internal/config"
)

type StaticHandler struct {
	Conf    *config.Config
	Handler http.Handler
}

func (h StaticHandler) Get(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = chi.URLParam(r, "*")
	h.Handler.ServeHTTP(w, r)
}
