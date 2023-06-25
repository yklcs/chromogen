package storage

import (
	"io"
	"net/http"
)

type Storage interface {
	Upload(r io.Reader, path string) string
	Url(path string)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
