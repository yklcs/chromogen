package photo

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
	"github.com/yklcs/wasmimg/mozjpeg"
)

func (p *Photo) ResizeAndCompress(longSideSize int, quality int) (int, int, error) {
	r, err := NewReader(*p)
	if err != nil {
		return 0, 0, err
	}
	var buf bytes.Buffer
	w, h, err := ResizeAndCompressStd(r, &buf, longSideSize, quality)
	if err != nil {
		return 0, 0, err
	}

	p.Close()
	p.Open()
	buf.WriteTo(p)

	return w, h, nil
}

func ResizeAndCompressStd(r io.Reader, w io.Writer, longSideSize int, quality int) (int, int, error) {
	img, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return 0, 0, err
	}

	if img.Bounds().Dx() > longSideSize {
		img = imaging.Resize(img, longSideSize, 0, imaging.Lanczos)
	}
	if img.Bounds().Dy() > longSideSize {
		img = imaging.Resize(img, 0, longSideSize, imaging.Lanczos)
	}

	imaging.Encode(w, img, imaging.JPEG, imaging.JPEGQuality(quality))
	return img.Bounds().Dx(), img.Bounds().Dy(), nil
}

func ResizeAndCompress(r io.Reader, w io.Writer, longSideSize int, quality int) (int, int, error) {
	rgb, width, height, err := mozjpeg.Decode(r)
	if err != nil {
		return 0, 0, err
	}
	img := rgbToImage(rgb, image.Rect(0, 0, width, height))

	var imgResized image.Image = img
	if img.Bounds().Dx() > longSideSize {
		imgResized = imaging.Resize(img, longSideSize, 0, imaging.Lanczos)
	}
	if img.Bounds().Dy() > longSideSize {
		imgResized = imaging.Resize(img, 0, longSideSize, imaging.Lanczos)
	}

	encoded, err := mozjpeg.Encode(bytes.NewReader(imageToRGB(imgResized)),
		imgResized.Bounds().Dx(),
		imgResized.Bounds().Dy(),
		quality,
	)
	if err != nil {
		return 0, 0, err
	}

	_, err = w.Write(encoded)
	return imgResized.Bounds().Dx(), imgResized.Bounds().Dy(), err
}

func ToJPEG(r io.Reader, w io.Writer, quality int) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	err = jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	return err
}

func rgbToImage(rgb []byte, bounds image.Rectangle) image.Image {
	img := image.NewRGBA(bounds)
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			rgbaIndex := (y*bounds.Dx() + x) * 4
			rgbIndex := (y*bounds.Dx() + x) * 3

			copy(img.Pix[rgbaIndex:rgbaIndex+3], rgb[rgbIndex:rgbIndex+3])
			img.Pix[rgbaIndex+3] = 255
		}
	}
	return img
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
			copy(rgb[rgbIndex:rgbIndex+3], pix)
		}
	}

	return rgb
}
