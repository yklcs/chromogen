//go:build !windows
// +build !windows

package photo

import (
	"image"
	"image/draw"
	"io"

	"github.com/disintegration/imaging"
	"github.com/yklcs/cram"
)

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
