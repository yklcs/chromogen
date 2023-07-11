package server

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/render"
)

func Auth(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")

			var token string
			if auth == "" {
				token = r.FormValue("token")
			} else {
				var ok bool
				token, ok = strings.CutPrefix(auth, "Bearer:")
				if !ok {
					http.Error(w, "malformed authorization", http.StatusUnauthorized)
					return
				}
				token = strings.TrimSpace(token)
			}

			if token == key {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "incorrect bearer token", http.StatusForbidden)
				return
			}
		})
	}
}

func AuthPage(key string, conf *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.FormValue("token")

			if token == key {
				r.Method = http.MethodGet
				chi.RouteContext(r.Context()).RouteMethod = http.MethodGet
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
				render.RenderAuth(w, conf)
				return
			}
		})
	}
}
