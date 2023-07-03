package render

import (
	"html/template"
	"io"
	"io/fs"
	"path"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photo"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/web"
)

type IndexedPhoto struct {
	photo.Photo
	Index int
}

func IndexPhoto(img photo.Photo, index int) IndexedPhoto {
	return IndexedPhoto{
		Photo: img,
		Index: index,
	}
}

func NewRootResolver(conf *config.Config) func(p ...string) string {
	return func(p ...string) string {
		p = append([]string{"/", conf.Root}, p...)
		return path.Join(p...)
	}
}

func RenderIndex(w io.Writer, ps *photos.Photos, conf *config.Config) error {
	type IndexTemplateData struct {
		Photos *photos.Photos
		*config.Config
	}

	templates, _ := fs.Sub(web.Content, "templates")
	tmpl := template.New("")
	tmpl = tmpl.Funcs(template.FuncMap{
		"IndexPhoto":   IndexPhoto,
		"RootResolver": NewRootResolver(conf),
	})
	tmpl, err := tmpl.ParseFS(templates, "index.tmpl", "head.tmpl", "thumb.tmpl")
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, "index", IndexTemplateData{ps, conf})
	if err != nil {
		return err
	}

	return nil
}

func RenderPhoto(w io.Writer, img *photo.Photo, conf *config.Config) error {
	type PhotoTemplateData struct {
		Photo  *photo.Photo
		Config *config.Config
	}

	templates, _ := fs.Sub(web.Content, "templates")
	tmpl := template.New("")
	tmpl = tmpl.Funcs(template.FuncMap{"RootResolver": NewRootResolver(conf)})
	tmpl, err := tmpl.ParseFS(templates, "photo.tmpl", "head.tmpl")
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, "photo", PhotoTemplateData{img, conf})
	if err != nil {
		return err
	}

	return nil
}

func RenderPanchro(w io.Writer, ps *photos.Photos, conf *config.Config) error {
	type PanchroTemplateData struct {
		Photos *photos.Photos
		*config.Config
	}

	templates, _ := fs.Sub(web.Content, "templates")
	tmpl := template.New("")
	tmpl = tmpl.Funcs(template.FuncMap{
		"RootResolver": NewRootResolver(conf),
	})
	tmpl, err := tmpl.ParseFS(templates, "panchro.tmpl", "head.tmpl", "thumb.tmpl")
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, "panchro", PanchroTemplateData{ps, conf})
	if err != nil {
		return err
	}

	return nil
}
