package server

import (
	"context"
	"net/http"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
)

type Server struct {
	Mux    *http.ServeMux
	Bucket *blob.Bucket
}

func NewServer(dir string) (*Server, error) {
	bucket, err := blob.OpenBucket(context.Background(), dir)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		Mux:    http.NewServeMux(),
		Bucket: bucket,
	}

	srv.InitMux()
	return srv, nil
}
