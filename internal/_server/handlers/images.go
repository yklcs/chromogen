package handlers

import (
	"net/http"
	"strings"

	"gocloud.dev/blob"
	"golang.org/x/crypto/scrypt"
)

type ImagesHandler struct {
	Bucket *blob.Bucket
}

func (h ImagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20+1024)
	// contentType := r.Header.Get("Content-Type")

	if r.Method != "POST" {
		http.Error(w, "unsupported method", http.StatusBadRequest)
		return
	}

	scrypt.Key([]byte("pw"), []byte("wow"), 32768, 8, 1, 32)

	err := r.ParseMultipartForm(32 << 20) // 32 MiB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, fileheader, _ := r.FormFile("file")
	defer file.Close()

	var contentType string = fileheader.Header.Get("Content-Type")
	if contentType == "" {
		sample := make([]byte, 512)
		_, err = file.Read(sample)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contentType = http.DetectContentType(sample)
	}
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
