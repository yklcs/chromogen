package photos

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
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

func (ps *Photos) LoadFiles(in []string, store storage.Storage) error {
	in, err := flattenPhotoPaths(in, []string{".jpeg", ".jpg", ".png"})
	if err != nil {
		return err
	}

	slices.Sort(in)
	for i, j := 0, len(in)-1; i < j; i, j = i+1, j-1 {
		in[i], in[j] = in[j], in[i]
	}

	jobs := make(chan string, len(in))
	results := make(chan *Photo, len(in))
	workers := 8
	bar := progressbar.Default(int64(len(in)), "LoadFiles")

	for w := 0; w < workers; w++ {
		go func(jobs <-chan string, results chan<- *Photo) {
			for j := range jobs {
				f, err := os.Open(j)
				if err != nil {
					log.Println(err)
				}
				p, err := NewPhoto(f, store)
				if err != nil {
					log.Println(err)
				}

				if _, err := ps.Get(p.ID); err != sql.ErrNoRows {
					p = nil
				}

				f.Close()
				bar.Add(1)
				results <- p
			}
		}(jobs, results)
	}
	for j := 0; j < len(in); j++ {
		jobs <- in[j]
	}
	close(jobs)
	for r := 0; r < len(in); r++ {
		ps.Set(<-results)
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
