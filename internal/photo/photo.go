package photo

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"html/template"
	"image"
	"io"
	"log"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/yklcs/panchro/internal/utils"
)

type Format string

const (
	Png  Format = ".png"
	Jpeg Format = ".jpg"
	Webp Format = ".webp"
)

type Photo struct {
	ID string `json:"id"`

	Title       string
	Description string
	Tags        []string

	URL  string `json:"url"`
	Path string `json:"path"`

	SourcePath string `json:"-"`

	Format Format `json:"format"`
	Hash   []byte `json:"-"`

	Exif *Exif `json:"exif"`

	PlaceholderURI template.URL `json:"-"`
	Width          int          `json:"width"`
	Height         int          `json:"height"`

	buffer *bytes.Buffer
}

type Exif struct {
	DateTime     time.Time `json:"datetime"`
	MakeModel    string    `json:"makemodel"`
	ShutterSpeed string    `json:"shutterspeed"`
	FNumber      string    `json:"fnumber"`
	ISO          string    `json:"iso"`
}

func NewPhoto(filepath string) (Photo, error) {
	var format Format

	ext := strings.ToLower(path.Ext(filepath))
	switch ext {
	case ".png":
		format = Png
	case ".jpg":
		format = Jpeg
	case ".jpeg":
		format = Jpeg
	case ".webp":
		format = Webp
	default:
		return Photo{}, errors.New("invalid format: " + ext)
	}

	return Photo{
		Format:     format,
		SourcePath: filepath,
	}, nil
}

func (p *Photo) ProcessMeta() error {
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		r, _ := NewReader(*p)

		hash := sha256.New()
		io.Copy(hash, r)
		p.Hash = hash.Sum(nil)

		p.ID = utils.Base58Encode(p.Hash)[:6]
		p.Path = p.ID + string(p.Format)
	}()

	go func() {
		defer wg.Done()
		r, _ := NewReader(*p)

		x, err := exif.Decode(r)
		if err != nil {
			log.Fatalln(err)
		}
		p.Exif = processExif(x)
	}()

	go func() {
		defer wg.Done()
		r, _ := NewReader(*p)

		img, _, err := image.DecodeConfig(r)
		if err != nil {
			log.Fatalln(err)
		}
		p.Width = img.Width
		p.Height = img.Height
	}()

	go func() {
		defer wg.Done()
		r, _ := NewReader(*p)

		placeholder := generatePlaceholderURI(r)
		p.PlaceholderURI = template.URL(placeholder)
	}()

	wg.Wait()
	return nil
}
