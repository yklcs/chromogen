package server

import (
	"net/http"
	"strings"
)

func Auth(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// auth := r.Header.Get("Authorization")
			auth := r.FormValue("token")
			// token, ok := strings.CutPrefix(auth, "Bearer:")
			// if !ok {
			// http.Error(w, "malformed authorization", http.StatusBadRequest)
			// return
			// }
			token := auth
			token = strings.TrimSpace(token)

			if token == key {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "incorrect bearer token", http.StatusForbidden)
				return
			}
		})
	}
}
