package photos

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/draw"
	"io"
	"log"
	"os"
	"path"
	"sync"

	"github.com/schollz/progressbar/v3"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
)

func generatePlaceholderURI(r io.Reader) string {
	var b bytes.Buffer
	ResizeAndCompress(r, &b, 12, 75)

	enc := base64.StdEncoding.EncodeToString(b.Bytes())
	return "data:image/jpeg;base64," + enc
}

func (ps *Photos) Upload(bucket *blob.Bucket, key string, p Photo) error {
	ctx := context.Background()

	w, err := bucket.NewWriter(ctx, "e", &blob.WriterOptions{})
	if err != nil {
		return err
	}
	defer w.Close()

	file, err := os.ReadFile(p.SourcePath)
	if err != nil {
		return err
	}

	_, err = w.Write(file)
	return err
}

func downloadPhoto(fpath string, r io.Reader) error {
	err := os.MkdirAll(path.Dir(fpath), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	return nil
}

func ListBucket(bucket *blob.Bucket) ([]string, error) {
	ctx := context.Background()

	var keys []string
	iter := bucket.List(nil)
	for {
		obj, err := iter.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		keys = append(keys, obj.Key)
	}

	return keys, nil
}

func (ps *Photos) Write(bucket *blob.Bucket) error {
	return nil
}

// Read photos, writing metadata into memory
func (ps *Photos) Read(bucket *blob.Bucket) error {
	ctx := context.Background()

	keys, err := ListBucket(bucket)
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(keys)), "download")
	var wg sync.WaitGroup
	wg.Add(len(keys))

	for _, key := range keys {
		go func(key string) {
			defer wg.Done()
			defer bar.Add(1)

			r, err := bucket.NewReader(ctx, key, nil)
			if err != nil {
				log.Fatalln(err)
			}
			defer r.Close()

			p, err := NewPhoto(key, ps.Bucket)
			if err != nil {
				log.Fatalln(err)
			}

			_, err = p.ReadFrom(r)
			if err != nil {
				log.Fatalln(err)
			}

			ps.Add(p)
		}(key)
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
