package render

import (
	"html/template"
	"io"
	"io/fs"
	"path"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/web"
)

type IndexedPhoto struct {
	photos.Photo
	Index int
}

func IndexPhoto(img photos.Photo, index int) IndexedPhoto {
	return IndexedPhoto{
		Photo: img,
		Index: index,
	}
}

func NewRootResolver(conf *config.Config) func(p string) string {
	return func(p string) string {
		return path.Join("/", conf.Root, p)
	}
}

func RenderIndex(w io.Writer, ps photos.Photos, conf *config.Config) error {
	type IndexTemplateData struct {
		Photos photos.Photos
		config.Config
	}

	templates, _ := fs.Sub(web.Content, "templates")
	tmpl := template.New("")
	tmpl = tmpl.Funcs(template.FuncMap{
		"IndexPhoto":   IndexPhoto,
		"RootResolver": NewRootResolver(conf),
	})
	tmpl, err := tmpl.ParseFS(templates, "index.tmpl", "thumb.tmpl")
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, "index", IndexTemplateData{ps, *conf})
	if err != nil {
		return err
	}

	return nil
}

func RenderPhoto(w io.Writer, img photos.Photo, conf *config.Config) error {
	type PhotoTemplateData struct {
		Photo  photos.Photo
		Config config.Config
	}

	templates, _ := fs.Sub(web.Content, "templates")
	tmpl := template.New("")
	tmpl = tmpl.Funcs(template.FuncMap{"RootResolver": NewRootResolver(conf)})
	tmpl, err := tmpl.ParseFS(templates, "photo.tmpl")
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, "photo", PhotoTemplateData{img, *conf})
	if err != nil {
		return err
	}

	return nil
}
