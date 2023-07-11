package render

import (
	"html/template"
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photo"
	"github.com/yklcs/panchro/internal/photos"
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

func RenderNewline(s string) template.HTML {
	return template.HTML(strings.ReplaceAll(s, "\n", "<br/>"))
}

func RenderIndex(w io.Writer, ps *photos.Photos, conf *config.Config) error {
	type IndexTemplateData struct {
		Photos *photos.Photos
		*config.Config
	}

	theme := config.LoadTheme(conf)
	templatesFS, _ := fs.Sub(theme, "templates")

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"IndexPhoto":    IndexPhoto,
		"RootResolver":  NewRootResolver(conf),
		"RenderNewline": RenderNewline,
	}).ParseFS(templatesFS, "index.tmpl", "head.tmpl", "thumb.tmpl")
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

	theme := config.LoadTheme(conf)
	templatesFS, _ := fs.Sub(theme, "templates")

	tmpl, err := template.New("").Funcs(template.FuncMap{"RootResolver": NewRootResolver(conf)}).ParseFS(templatesFS, "photo.tmpl", "head.tmpl")
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

	theme := config.LoadTheme(conf)
	templatesFS, _ := fs.Sub(theme, "templates")

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"RootResolver": NewRootResolver(conf),
	}).ParseFS(templatesFS, "panchro.tmpl", "head.tmpl", "thumb.tmpl")
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, "panchro", PanchroTemplateData{ps, conf})
	if err != nil {
		return err
	}

	return nil
}

func RenderAuth(w io.Writer, conf *config.Config) error {
	type AuthTemplateData struct {
		*config.Config
	}

	theme := config.LoadTheme(conf)
	templatesFS, _ := fs.Sub(theme, "templates")

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"RootResolver": NewRootResolver(conf),
	}).ParseFS(templatesFS, "auth.tmpl", "head.tmpl")
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(w, "auth", AuthTemplateData{conf})
	if err != nil {
		return err
	}

	return nil
}
