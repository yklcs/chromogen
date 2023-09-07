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
	ID string `json:"id"`

	Title       string
	Description string
	Tags        []string

	URL  string `json:"url"`
	Path string `json:"path"`

	SourcePath string `json:"-"`

	Format string `json:"format"`
	Hash   []byte `json:"hash"`

	Exif *Exif `json:"exif"`

	PlaceholderURI template.URL `json:"placeholder_uri"`
	Width          int          `json:"width"`
	Height         int          `json:"height"`

	data []byte
}

type Exif struct {
	DateTime        time.Time `json:"datetime"`
	MakeModel       string    `json:"makemodel"`
	ShutterSpeed    string    `json:"shutterspeed"`
	FNumber         string    `json:"fnumber"`
	ISO             string    `json:"iso"`
	LensMakeModel   string    `json:"lens_makemodel"`
	FocalLength     string    `json:"focallength"`
	SubjectDistance string    `json:"subjectdistance"`
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
	}
	p.PlaceholderURI = template.URL(generatePlaceholderURI(bytes.NewReader(p.data)))

	return &p, nil
}
