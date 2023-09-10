package photos

import (
	"bytes"
	"crypto/sha256"
	"html/template"
	"image"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/yklcs/panchro/internal/utils"
)

type Photo struct {
	ID string

	Title       string
	Description string
	Tags        []string

	URL  string
	Path string

	Format string
	Hash   []byte

	Exif *Exif

	PlaceholderURI template.URL
	Width          int
	Height         int

	data []byte
}

type Exif struct {
	DateTime        time.Time
	MakeModel       string
	ShutterSpeed    string
	FNumber         string
	ISO             string
	LensMakeModel   string
	FocalLength     string
	SubjectDistance string
}

func NewPhotoFromFile(filepath string) (*Photo, error) {
	var p Photo
	var err error

	p.data, err = os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(p.data))
	if err != nil {
		return nil, err
	}
	p.Width = cfg.Width
	p.Height = cfg.Height
	p.Format = format

	hash := sha256.Sum256(p.data)
	p.Hash = hash[:]
	p.ID = utils.Base58Encode(p.Hash)[:6]
	p.Path = p.ID + "." + p.Format

	x, err := exif.Decode(bytes.NewReader(p.data))
	if err == nil {
		p.Exif = processExif(x)
	} else {
		p.Exif = &Exif{}
	}
	p.PlaceholderURI = template.URL(
		generatePlaceholderURI(bytes.NewReader(p.data)))

	return &p, nil
}
