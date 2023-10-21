package build

import (
	"database/sql"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/storage"
	"golang.org/x/exp/slices"
)

type StaticSiteGenerator struct {
	outpath string
	conf    *config.Config
	photos  *photos.Photos
}

func NewStaticSiteGenerator(outpath, confpath string) (*StaticSiteGenerator, error) {
	conf, err := config.ReadConfig(confpath)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(outpath, 0755)
	if err != nil {
		return nil, err
	}

	dbpath := path.Join(outpath, "panchro.db")
	db, err := sql.Open("sqlite", dbpath)
	if err != nil {
		return nil, err
	}

	return &StaticSiteGenerator{
		outpath: outpath,
		conf:    conf,
		photos:  &photos.Photos{DB: db},
	}, nil
}

func (s *StaticSiteGenerator) Build(inpaths []string) error {
	defer s.photos.DB.Close()
	err := s.photos.Init()
	if err != nil {
		return err
	}

	inpaths, err = flattenPhotoPaths(inpaths, []string{".jpg", ".jpeg", ".png"})
	if err != nil {
		return err
	}

	store, _ := storage.NewLocalStorage(s.outpath)
	err = s.photos.LoadFiles(inpaths, store)
	if err != nil {
		return err
	}

	theme, err := config.NewTheme(s.conf)
	if err != nil {
		return err
	}

	indexHTML, err := os.Create(path.Join(s.outpath, "index.html"))
	if err != nil {
		return err
	}
	defer indexHTML.Close()

	err = theme.Render(indexHTML, "index",
		config.ThemeData{Photos: s.photos, Config: s.conf})
	if err != nil {
		return err
	}

	staticDir := path.Join(s.outpath, s.conf.StaticDir)
	err = os.MkdirAll(staticDir, 0755)
	if err != nil {
		return err
	}
	err = theme.WriteStatic(staticDir)
	if err != nil {
		return err
	}

	for _, id := range s.photos.IDs() {
		dir := path.Join(s.outpath, id)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		imageHTML, err := os.Create(path.Join(dir, "index.html"))
		if err != nil {
			return err
		}
		defer imageHTML.Close()

		p, _ := s.photos.Get(id)
		err = theme.Render(imageHTML, "photo", config.ThemeData{Photo: p, Config: s.conf})
		if err != nil {
			return err
		}
	}

	return nil
}

func flattenPhotoPaths(dirs []string, exts []string) ([]string, error) {
	var matched []string
	for _, dir := range dirs {
		err := filepath.WalkDir(dir, func(s string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if d.IsDir() {
				return nil
			}
			if slices.Contains(exts, filepath.Ext(d.Name())) {
				matched = append(matched, s)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return matched, nil
}
