package photos

import (
	"bytes"
	"crypto/sha256"
	"html/template"
	"image"
	"io"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/yklcs/chromogen/internal/utils"
	"github.com/yklcs/chromogen/storage"
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

func NewPhoto(r io.Reader, store storage.Storage) (*Photo, error) {
	var p Photo
	var buf bytes.Buffer

	buf.ReadFrom(r)

	cfg, format, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	p.Width = cfg.Width
	p.Height = cfg.Height
	p.Format = format

	hash := sha256.Sum256(buf.Bytes())
	p.Hash = hash[:]
	p.ID = utils.Base58Encode(p.Hash)[:6]
	p.Path = p.ID + "." + p.Format

	x, err := exif.Decode(bytes.NewReader(buf.Bytes()))
	if err == nil {
		p.Exif = processExif(x)
	} else {
		p.Exif = &Exif{}
	}
	p.PlaceholderURI = template.URL(
		generatePlaceholderURI(bytes.NewReader(buf.Bytes())))

	url, err := store.Upload(bytes.NewReader(buf.Bytes()), p.Path)
	if err != nil {
		return nil, err
	}
	p.URL = url

	return &p, nil
}
