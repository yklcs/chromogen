package storage

import (
	"io"
	"net/http"
)

type Storage interface {
	Upload(r io.Reader, fpath string) (string, error)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Backend() string
}
