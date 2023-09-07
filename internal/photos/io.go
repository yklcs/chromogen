package photos

import (
	"bytes"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
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

func (ps *Photos) Upload(store storage.Storage) error {
	ids := ps.IDs()
	for _, id := range ids {
		p, err := ps.Get(id)
		if err != nil {
			return err
		}
		url, err := store.Upload(bytes.NewReader(p.data), p.Path)
		if err != nil {
			return err
		}
		p.URL = url
		ps.Set(p)
	}
	return nil
}

func (ps *Photos) LoadFS(in string) error {
	fpaths, err := MatchExts(in, []string{".jpeg", ".jpg"})
	if err != nil {
		return err
	}
	slices.Sort(fpaths)
	for i, j := 0, len(fpaths)-1; i < j; i, j = i+1, j-1 {
		fpaths[i], fpaths[j] = fpaths[j], fpaths[i]
	}

	jobs := make(chan string, len(fpaths))
	results := make(chan *Photo, len(fpaths))
	workers := 8
	bar := progressbar.Default(int64(len(fpaths)), "LoadFS")

	for w := 0; w < workers; w++ {
		go func(jobs <-chan string, results chan<- *Photo) {
			for j := range jobs {
				p, err := NewPhotoFromFile(j)
				if err != nil {
					log.Println(err)
				}
				bar.Add(1)
				results <- p
			}
		}(jobs, results)
	}
	for j := 0; j < len(fpaths); j++ {
		jobs <- fpaths[j]
	}
	close(jobs)
	for r := 0; r < len(fpaths); r++ {
		ps.Set(<-results)
	}

	return nil
}

// ProcessFS reads and processes photo data from the filesystem
// func (ps *Photos) ProcessFS(in, out string, compress bool, longSideSize, quality int) error {
// 	oldps := &Photos{}
// 	oldpsExists, err := oldps.Load(path.Join(out, "panchro.db"))

// 	store, err := storage.NewLocalStorage(out)
// 	if err != nil {
// 		return err
// 	}

// 	fpaths, err := MatchExts(in, []string{".jpeg", ".jpg"})
// 	if err != nil {
// 		return err
// 	}
// 	slices.Sort(fpaths)

// 	bar := progressbar.Default(int64(len(fpaths)), "process")

// 	var wg sync.WaitGroup
// 	wg.Add(len(fpaths))

// 	for _, fpath := range fpaths {
// 		fin, err := os.Open(fpath)
// 		if err != nil {
// 			return err
// 		}
// 		defer fin.Close()

// 		relpath, err := filepath.Rel(in, fpath)
// 		if err != nil {
// 			return err
// 		}

// 		hash := sha256.New()
// 		io.Copy(hash, fin)
// 		phash := hash.Sum(nil)
// 		pid := utils.Base58Encode(phash)[:6]
// 		fin.Seek(0, 0)

// 		if oldpsExists {
// 			if oldp, err := oldps.Get(pid); err == nil {
// 				// photo already exists in db, use that
// 				ps.Add(oldp)
// 				wg.Done()
// 				bar.Add(1)
// 				continue
// 			}
// 		}

// 		p, err := photo.NewPhoto(relpath)
// 		if err != nil {
// 			return err
// 		}
// 		p.Open()
// 		p.ReadFrom(fin)

// 		err = p.ProcessMeta()
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		w, h := p.Width, p.Height
// 		var outbuf bytes.Buffer
// 		r, _ := photo.NewReader(p)
// 		if compress {
// 			w, h, _ = p.ResizeAndCompress(longSideSize, quality)
// 		} else {
// 			photo.ToJPEG(r, &outbuf, quality)
// 		}
// 		p.Width = w
// 		p.Height = h

// 		r, _ = photo.NewReader(p)
// 		purl, err := store.Upload(r, p.Path)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		p.URL = purl

// 		p.Close()
// 		ps.Add(p)

// 		wg.Done()
// 		bar.Add(1)
// 	}

// 	wg.Wait()
// 	return nil
// }
