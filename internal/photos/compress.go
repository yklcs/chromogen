//go:build !windows
// +build !windows

package photos

import (
	"image"
	"io"
	"log"
	"os"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/schollz/progressbar/v3"
	"github.com/yklcs/cram"
)

func compressPhoto(p *Photo, src string, dst string, longSideSize int, quality int) {
	srcf, err := os.Open(src)
	if err != nil {
		log.Fatalln(err)
	}
	defer srcf.Close()

	dstf, err := os.Create(dst)
	if err != nil {
		log.Fatalln(err)
	}
	defer dstf.Close()

	w, h, err := ResizeAndCompress(srcf, dstf, longSideSize, quality)
	if err != nil {
		log.Fatalln(err)
	}

	p.Width = w
	p.Height = h
}

func (ps *Photos) Compress() error {
	num_images := ps.Len()
	bar := progressbar.Default(int64(num_images), "compress")

	var wg sync.WaitGroup
	wg.Add(ps.Len())
	for i, p := range ps.Slice() {
		go func(i int, p Photo) {
			defer wg.Done()
			defer bar.Add(1)
			compressPhoto(ps.Get(i), p.OriginalFilePath, p.FilePath, 2048, 80)
		}(i, p)
	}

	wg.Wait()

	return nil
}

func ResizeAndCompress(r io.Reader, w io.Writer, longSideSize int, quality int) (int, int, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return 0, 0, err
	}

	var imgResized image.Image = img
	if img.Bounds().Dx() > longSideSize {
		imgResized = imaging.Resize(img, longSideSize, 0, imaging.Lanczos)
	}
	if img.Bounds().Dy() > longSideSize {
		imgResized = imaging.Resize(img, 0, longSideSize, imaging.Lanczos)
	}

	compressed, err := cram.MozJPEG(imageToRGB(imgResized),
		imgResized.Bounds().Dx(), imgResized.Bounds().Dy(), quality)
	if err != nil {
		return 0, 0, err
	}

	_, err = w.Write(compressed)
	return imgResized.Bounds().Dx(), imgResized.Bounds().Dy(), err
}
