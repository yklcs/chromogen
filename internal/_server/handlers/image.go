package handlers

import (
	"net/http"
	"strings"

	"gocloud.dev/blob"
)

type ImageHandler struct {
	Bucket *blob.Bucket
}

type ImageRequest struct {
	Key string `json:"key"`
}

func (h ImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// var req ImageRequest
	// err := json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	switch r.Method {
	case "GET":
		path, _ := strings.CutPrefix(r.URL.Path, "/images/")
		h := http.FileServer(http.Dir("/"))
		h.ServeHTTP(w, r)
	case "POST":
		h.Bucket.NewWriter(r.Context(), "e", nil)
		// h.Store.Put(req.Key, []byte("world"))
	}
	// w.Write([]byte(""))
}
