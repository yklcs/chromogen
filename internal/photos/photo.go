package photos

import (
	"bytes"
	"context"
	"crypto/sha256"
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
	"gocloud.dev/blob"
)

type Format string

const (
	Png  Format = ".png"
	Jpeg        = ".jpg"
	Webp        = ".webp"
)

type Photo struct {
	ID string

	Path       string
	SourcePath string

	Format Format
	Hash   []byte

	Exif *Exif

	PlaceholderURI template.URL
	Width          int
	Height         int

	bucket *blob.Bucket
}

type Exif struct {
	DateTime     time.Time
	MakeModel    string
	ShutterSpeed string
	FNumber      string
	ISO          string
}

func NewPhoto(filepath string, bucket *blob.Bucket) (Photo, error) {
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
		bucket:     bucket,
	}, nil
}

func (p *Photo) Read(b []byte) (int, error) {
	r, err := p.bucket.NewReader(context.Background(), p.Path, nil)
	if err != nil {
		return len(b), err
	}
	defer r.Close()

	n, err := r.Read(b)
	return n, err
}

func (p *Photo) ReadFrom(r io.Reader) (int, error) {
	buf := bytes.Buffer{}
	buf.ReadFrom(r)

	if len(p.Hash) == 0 {
		hash := sha256.New()

		hashN, err := hash.Write(buf.Bytes())
		if err != nil {
			return hashN, err
		}
		p.Hash = hash.Sum(nil)
		p.ID = utils.Base58Encode(p.Hash)[:8]
		p.Path = p.ID + string(p.Format)
	}

	if p.Exif == nil {
		x, err := exif.Decode(bytes.NewReader(buf.Bytes()))
		if err != nil {
			return buf.Len(), err
		}
		p.Exif = ProcessExif(x)
	}

	img, _, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return buf.Len(), err
	}
	p.Width = img.Width
	p.Height = img.Height

	if p.PlaceholderURI == "" {
		placeholder := generatePlaceholderURI(bytes.NewReader(buf.Bytes()))
		p.PlaceholderURI = template.URL(placeholder)
	}

	w, err := p.bucket.NewWriter(context.Background(), p.Path, nil)
	if err != nil {
		return buf.Len(), err
	}
	defer w.Close()

	io.Copy(w, bytes.NewReader(buf.Bytes()))
	// _, err = w.Write(b)
	if err != nil {
		return buf.Len(), err
	}

	return buf.Len(), nil
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
