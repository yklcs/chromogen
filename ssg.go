package panchro

import (
	"database/sql"
	"io/fs"
	"os"
	"path"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/render"
)

type StaticSiteGenerator struct {
	inpath  string
	outpath string
	conf    *config.Config
	photos  *photos.Photos
}

func NewStaticSiteGenerator(inpath, outpath, confpath string) (*StaticSiteGenerator, error) {
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
		inpath:  inpath,
		outpath: outpath,
		conf:    conf,
		photos:  &photos.Photos{DB: db},
	}, nil
}

func (s *StaticSiteGenerator) Build() error {
	defer s.photos.DB.Close()
	err := s.photos.Init()
	if err != nil {
		return err
	}

	err = s.photos.ProcessFS(s.inpath, s.outpath, true, 2048, 75)
	if err != nil {
		return err
	}

	indexHTML, err := os.Create(path.Join(s.outpath, "index.html"))
	if err != nil {
		return err
	}
	defer indexHTML.Close()

	err = render.RenderIndex(indexHTML, s.photos, s.conf)
	if err != nil {
		return err
	}

	themeFS := config.LoadTheme(s.conf)
	staticFS, _ := fs.Sub(themeFS, "static")
	staticDir := path.Join(s.outpath, s.conf.StaticDir)
	err = os.MkdirAll(staticDir, 0755)
	if err != nil {
		return err
	}
	err = render.CopyFS(staticFS, staticDir)
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
		err = render.RenderPhoto(imageHTML, &p, s.conf)
		if err != nil {
			return err
		}
	}

	return nil
}
