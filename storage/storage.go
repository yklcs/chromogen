package storage

import (
	"io"
	"net/http"
)

type Backend string

const (
	S3Backend    Backend = "s3"
	LocalBackend Backend = "local"
)

type Storage interface {
	Upload(r io.Reader, fpath string) (string, error)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Backend() string
}

type Reader struct {
	i int64
}
