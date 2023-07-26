package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yklcs/panchro/internal/photo"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/storage"
)

type PhotosHandler struct {
	Photos *photos.Photos
	Store  storage.Storage
}

func (h PhotosHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.Photos)
	w.Header().Set("Content-Type", "application/json")
}

func (h PhotosHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.Photos.Get(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(p)
	w.Header().Set("Content-Type", "application/json")
}

func (h PhotosHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.Photos.Delete(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not delete file", http.StatusInternalServerError)
		return
	}
}

func (h PhotosHandler) Post(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) // 32MB
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "could not upload file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	p, _ := photo.NewPhoto(handler.Filename)
	p.Open()
	p.ReadFrom(file)
	p.ProcessMeta()
	p.ResizeAndCompress(2048, 75)
	p.Upload(h.Store)
	h.Photos.Add(p)
	// w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
