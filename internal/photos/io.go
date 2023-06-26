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
	_ "gocloud.dev/blob/s3blob"
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
func (ps *Photos) ProcessFS(in, out string, longSideSize, quality int) error {
	store, err := storage.NewLocalStorage(out)
	ps.store = store
	if err != nil {
		return err
	}

	fpaths, err := MatchExts(in, []string{".jpeg", ".jpg"})
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(fpaths)), "read")

	var wg sync.WaitGroup
	wg.Add(len(fpaths))

	for _, fpath := range fpaths {
		fin, err := os.Open(fpath)
		if err != nil {
			return err
		}
		defer fin.Close()
		var buf bytes.Buffer
		buf.ReadFrom(fin)

		relpath, err := filepath.Rel(in, fpath)
		if err != nil {
			return err
		}

		p, err := photo.NewPhoto(relpath)
		if err != nil {
			return err
		}

		metaDone := make(chan bool, 1)
		go func() {
			err = p.ProcessMeta(bytes.NewReader(buf.Bytes()))
			if err != nil {
				log.Fatalln(err)
			}
			ps.Add(p)
			metaDone <- true
		}()

		go func() {
			defer wg.Done()
			defer bar.Add(1)

			var buf2 bytes.Buffer
			w, h, err := photo.ResizeAndCompress(bytes.NewReader(buf.Bytes()), &buf2, longSideSize, quality)
			if err != nil {
				log.Fatalln(err)
			}

			<-metaDone
			p.Width = w
			p.Height = h

			p.Upload(bytes.NewReader(buf2.Bytes()), store)
		}()
	}

	wg.Wait()
	return nil
}
