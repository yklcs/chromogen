package photos

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"github.com/yklcs/chromogen/internal/config"
	"github.com/yklcs/chromogen/storage"
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

func (ps *Photos) LoadFiles(in []string, store storage.Storage, conf *config.Config) error {
	in, err := flattenPhotoPaths(in, []string{".jpeg", ".jpg", ".png"})
	if err != nil {
		return err
	}

	slices.Sort(in)
	for i, j := 0, len(in)-1; i < j; i, j = i+1, j-1 {
		in[i], in[j] = in[j], in[i]
	}

	var deltaPaths []string
	for _, fpath := range in {
		f, err := os.Open(fpath)
		if err != nil {
			return err
		}

		id := PhotoId(f)

		if _, err := ps.Get(id); err == sql.ErrNoRows {
			deltaPaths = append(deltaPaths, fpath)
		}
	}

	jobs := make(chan string, len(deltaPaths))
	results := make(chan *Photo, len(deltaPaths))
	workers := 8
	bar := progressbar.Default(int64(len(deltaPaths)), "LoadFiles")

	for w := 0; w < workers; w++ {
		go func(jobs <-chan string, results chan<- *Photo) {
			for j := range jobs {
				f, err := os.Open(j)
				if err != nil {
					log.Println(err)
				}
				p, err := NewPhoto(f, store, conf.ThumbSize)
				if err != nil {
					log.Println(err)
				}
				f.Close()
				p.srcPath = j
				bar.Add(1)
				results <- p
			}
		}(jobs, results)
	}
	for j := 0; j < len(deltaPaths); j++ {
		jobs <- deltaPaths[j]
	}
	close(jobs)

	delta := make([]*Photo, len(deltaPaths))
	for i := 0; i < len(deltaPaths); i++ {
		delta[i] = <-results
	}

	slices.SortStableFunc(delta, func(a, b *Photo) bool {
		return a.srcPath < b.srcPath
	})

	for _, d := range delta {
		ps.Set(d)
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
