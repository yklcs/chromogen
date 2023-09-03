package photos

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/yklcs/panchro/internal/photo"
	"github.com/yklcs/panchro/storage"
	"golang.org/x/exp/slices"
)

func MatchExts(dir string, exts []string) ([]string, error) {
	var matched []string
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
	return matched, err
}

// ProcessFS reads and processes photo data from the filesystem
func (ps *Photos) ProcessFS(in, out string, compress bool, longSideSize, quality int) error {
	store, err := storage.NewLocalStorage(out)
	if err != nil {
		return err
	}

	fpaths, err := MatchExts(in, []string{".jpeg", ".jpg"})
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(fpaths)), "process")

	var wg sync.WaitGroup
	wg.Add(len(fpaths))

	for _, fpath := range fpaths {
		fin, err := os.Open(fpath)
		if err != nil {
			return err
		}
		defer fin.Close()

		relpath, err := filepath.Rel(in, fpath)
		if err != nil {
			return err
		}

		p, err := photo.NewPhoto(relpath)
		if err != nil {
			return err
		}
		p.Open()
		p.ReadFrom(fin)

		metaDone := make(chan bool, 1)
		func() {
			err = p.ProcessMeta()
			if err != nil {
				log.Fatalln(err)
			}
			ps.Add(p)
			metaDone <- true
		}()

		go func() {
			defer wg.Done()
			defer bar.Add(1)

			w, h := p.Width, p.Height
			var outbuf bytes.Buffer
			r, _ := photo.NewReader(p)
			if compress {
				p.ResizeAndCompress(longSideSize, quality)
			} else {
				photo.ToJPEG(r, &outbuf, quality)
			}
			<-metaDone

			p.Width = w
			p.Height = h

			r, _ = photo.NewReader(p)
			purl, err := store.Upload(r, p.Path)
			if err != nil {
				log.Fatalln(err)
			}
			p.URL = purl

			p.Close()
			ps.Add(p)
		}()
	}

	wg.Wait()
	return nil
}
