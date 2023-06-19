package handlers

import (
	"io/fs"
	"net/http"
	"regexp"
	"strings"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/render"
	"github.com/yklcs/panchro/web"
)

type IndexHandler struct {
	Photos photos.Photos
	Conf   *config.Config
}

func (h IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jpgPattern, _ := regexp.Compile("[[:alnum:]]*.jpg$")
	staticPattern, _ := regexp.Compile("[[:alnum:]]*.woff2|js|css$")

	staticFs, _ := fs.Sub(web.Content, "static")
	staticHandler := http.FileServer(http.FS(staticFs))

	if r.URL.Path == "/" {
		render.RenderIndex(w, h.Photos, h.Conf)
	} else if jpg := jpgPattern.FindString(r.URL.Path); jpg != "" {
		http.ServeFile(w, r, "dist"+"/"+jpg)
	} else if static := staticPattern.FindString(r.URL.Path); static != "" {
		staticHandler.ServeHTTP(w, r)
	} else {
		render.RenderPhoto(w, *h.Photos.Find(strings.Split(r.URL.Path, "/")[1]), h.Conf)
	}
}
