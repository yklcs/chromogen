package photos

import (
	"image"
	"image/draw"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/schollz/progressbar/v3"
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

// Read photo metadata into memory and photo data into FS
func (ps *Photos) ReadFS(dir string) error {
	ps.store, err = storage.NewLocalStorage("panchro")
	if err != nil {
		return err
	}

	fpaths, err := MatchExts(dir, []string{".jpeg", ".jpg"})
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(fpaths)), "download")
	var wg sync.WaitGroup
	wg.Add(len(fpaths))

	for _, fpath := range fpaths {
		go func(fpath string) {
			defer wg.Done()
			defer bar.Add(1)

			p, err := NewPhoto(fpath)
			if err != nil {
				log.Fatalln(err)
			}

			f, err := os.Open(fpath)
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()

			err = p.ProcessMeta(f)
			if err != nil {
				log.Fatalln(err)
			}

			p.Upload()

			ps.Add(p)
		}(fpath)
	}

	wg.Wait()
	return nil
}

func imageToRGB(img image.Image) []byte {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	rgb := make([]byte, width*height*3)
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rgbaIndex := (y*width + x) * 4
			rgbIndex := (y*width + x) * 3
			pix := rgba.Pix[rgbaIndex : rgbaIndex+4]
			rgb[rgbIndex] = pix[0]
			rgb[1] = pix[1]
			rgb[2] = pix[2]
			copy(rgb[rgbIndex:rgbIndex+3], pix)
		}
	}

	return rgb
}
