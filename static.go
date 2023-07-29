package panchro

import (
	"io/fs"
	"os"
	"path"

	"github.com/yklcs/panchro/internal/config"
	"github.com/yklcs/panchro/internal/photos"
	"github.com/yklcs/panchro/internal/render"
	bolt "go.etcd.io/bbolt"
)

type StaticSiteGenerator struct {
	in     string
	out    string
	conf   *config.Config
	photos *photos.Photos
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

	dbpath := path.Join(outpath, "tmp.db")
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &StaticSiteGenerator{
		in:     inpath,
		out:    outpath,
		conf:   conf,
		photos: &photos.Photos{DB: db},
	}, nil
}

func (s *StaticSiteGenerator) Build() error {
	defer s.photos.DB.Close()
	defer os.RemoveAll(s.photos.DB.Path())
	s.photos.Init()

	err := s.photos.ProcessFS(s.in, s.out, true, 2048, 75)
	if err != nil {
		return err
	}

	indexHTML, err := os.Create(path.Join(s.out, "index.html"))
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
	staticDir := path.Join(s.out, s.conf.StaticDir)
	err = os.MkdirAll(staticDir, 0755)
	if err != nil {
		return err
	}
	err = render.CopyFS(staticFS, staticDir)
	if err != nil {
		return err
	}

	for _, id := range s.photos.IDs() {
		dir := path.Join(s.out, id)
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
