package config

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/yklcs/chromogen/internal/photos"
	"github.com/yklcs/chromogen/theme"
)

func LoadTheme(conf *Config) fs.FS {
	var themeFS fs.FS
	if conf.Theme == "" {
		return theme.FS
	} else {
		themeFS = os.DirFS(conf.Theme)
	}
	return themeFS
}

type Theme struct {
	staticFS     fs.FS
	templatesFS  fs.FS
	templateRoot *template.Template
}

func NewTheme(conf *Config) (*Theme, error) {
	var themeFS fs.FS
	if conf.Theme == "" {
		themeFS = theme.FS
	} else {
		themeFS = os.DirFS(conf.Theme)
	}
	templatesFS, _ := fs.Sub(themeFS, "templates")
	staticFS, _ := fs.Sub(themeFS, "static")

	funcmap := template.FuncMap{
		"Map": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid map call")
			}
			data := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("map keys must be strings")
				}
				data[key] = values[i+1]
			}
			return data, nil
		},
		"RootResolver": func(p ...string) string {
			if !strings.HasPrefix(p[0], "https://") {
				p = append([]string{"/", conf.Root}, p...)
			}
			return path.Join(p...)
		},
	}

	tmpl, err := template.New("").Funcs(funcmap).ParseFS(templatesFS, "*.tmpl")
	if err != nil {
		return nil, err
	}

	return &Theme{
		staticFS:     staticFS,
		templatesFS:  templatesFS,
		templateRoot: tmpl,
	}, err
}

type ThemeData struct {
	Photos *photos.Photos
	Photo  *photos.Photo
	Config *Config
}

func (t *Theme) Render(w io.Writer, name string, data ThemeData) error {
	return t.templateRoot.ExecuteTemplate(w, name, data)
}

func (t *Theme) WriteStatic(dst string) error {
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	return fs.WalkDir(t.staticFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			err = os.MkdirAll(path.Join(dst, p), 0755)
			if err != nil {
				return err
			}
		} else {
			fsrc, err := t.staticFS.Open(p)
			if err != nil {
				return err
			}
			defer fsrc.Close()

			fdst, err := os.Create(path.Join(dst, p))
			if err != nil {
				return err
			}
			defer fdst.Close()

			io.Copy(fdst, fsrc)
		}

		return nil
	})
}
