package photos

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"path"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/yklcs/panchro/internal/utils"
	"github.com/yklcs/panchro/storage"
)

type Format string

const (
	Png  Format = ".png"
	Jpeg        = ".jpg"
	Webp        = ".webp"
)

type Photo struct {
	ID  string
	URL string

	Path       string
	SourcePath string

	Format Format
	Hash   []byte

	Exif *Exif

	PlaceholderURI template.URL
	Width          int
	Height         int
}

type Exif struct {
	DateTime     time.Time
	MakeModel    string
	ShutterSpeed string
	FNumber      string
	ISO          string
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

func (p *Photo) ProcessMeta(r io.Reader) error {
	hashReader, hashWriter := io.Pipe()
	exifReader, exifWriter := io.Pipe()
	imageReader, imageWriter := io.Pipe()
	placeholderReader, placeholderWriter := io.Pipe()
	w := io.MultiWriter(hashWriter, exifWriter, imageWriter, placeholderWriter)

	if len(p.Hash) == 0 {
		hash := sha256.New()

		_, err := io.Copy(hash, hashReader)
		if err != nil {
			return err
		}
		p.Hash = hash.Sum(nil)
		p.ID = utils.Base58Encode(p.Hash)[:8]
		p.Path = p.ID + string(p.Format)
	}

	if p.Exif == nil {
		x, err := exif.Decode(exifReader)
		if err != nil {
			return err
		}
		p.Exif = ProcessExif(x)
	}

	img, _, err := image.DecodeConfig(imageReader)
	if err != nil {
		return err
	}
	p.Width = img.Width
	p.Height = img.Height

	if p.PlaceholderURI == "" {
		placeholder := generatePlaceholderURI(placeholderReader)
		p.PlaceholderURI = template.URL(placeholder)
	}

	_, err = io.Copy(w, r)
	return err
}

func (p *Photo) Upload(r io.Reader, store storage.Storage) error {
	purl, err := store.Upload(r, p.Path)
	p.URL = purl
	return err
}

func (p *Photo) Compress(longSideSize int, quality int) {
	pr, pw := io.Pipe()
	p.ReadFrom(pr)
	w, h, err := ResizeAndCompress(p, pw, longSideSize, quality)

	p.Width = w
	p.Height = h

	if err != nil {
		log.Fatalln(err)
	}
}

func ProcessExif(x *exif.Exif) *Exif {
	ex := Exif{}

	mkTag, _ := x.Get(exif.Make)
	if mkTag != nil {
		ex.MakeModel, _ = mkTag.StringVal()
	}

	modelTag, _ := x.Get(exif.Model)
	if modelTag != nil {
		model, _ := modelTag.StringVal()
		ex.MakeModel += " " + model
	}

	ex.DateTime, _ = x.DateTime()

	ssTag, _ := x.Get(exif.ShutterSpeedValue)
	if ssTag != nil {
		ssApexRat, _ := ssTag.Rat(0)
		ssApex, _ := ssApexRat.Float64()
		ss := math.Pow(2, -ssApex)

		if ss > 1 {
			ex.ShutterSpeed = fmt.Sprintf("%.1fs", ss)
		} else {
			ex.ShutterSpeed = fmt.Sprintf("1/%ds", int(math.Round(1/ss)))
		}
	}

	fTag, _ := x.Get(exif.FNumber)
	if fTag != nil {
		fRat, _ := fTag.Rat(0)
		ex.FNumber = "ƒ/"
		if fRat.IsInt() {
			ex.FNumber += fRat.RatString()
		} else {
			f, _ := fRat.Float64()
			ex.FNumber += fmt.Sprintf("%.1f", f)
		}
	}

	isoTag, _ := x.Get(exif.ISOSpeedRatings)
	if isoTag != nil {
		iso, _ := isoTag.Int64(0)
		ex.ISO = "ISO " + fmt.Sprint(iso)
	}

	return &ex
}

func generatePlaceholderURI(r io.Reader) string {
	var b bytes.Buffer
	ResizeAndCompress(r, &b, 12, 75)

	enc := base64.StdEncoding.EncodeToString(b.Bytes())
	return "data:image/jpeg;base64," + enc
}
