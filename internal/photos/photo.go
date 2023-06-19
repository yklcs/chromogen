package photos

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"html/template"
	"image"
	"io"
	"math"
	"path"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type Format int

const (
	Png Format = iota
	Jpeg
	Webp
)

type Photo struct {
	ID string

	Path             string
	SourcePath       string
	FilePath         string
	OriginalFilePath string

	Format   Format
	DateTime time.Time
	Hash     []byte

	MakeModel    string
	ShutterSpeed string
	FNumber      string
	ISO          string

	PlaceholderURI template.URL
	Width          int
	Height         int
}

func NewPhoto(imgPath string, dir string, r io.Reader) (Photo, error) {
	var format Format

	ext := strings.ToLower(path.Ext(imgPath))
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
		return Photo{}, errors.New("invalid format")
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return Photo{}, err
	}

	hash := sha256.New()
	_, err = hash.Write(buf)
	if err != nil {
		return Photo{}, err
	}

	x, err := exif.Decode(bytes.NewReader(buf))
	if err != nil {
		return Photo{}, err
	}

	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return Photo{}, err
	}

	hashstr := base58Encode(hash.Sum(nil))
	id := hashstr[:8]

	fpath := id + ext

	p := ProcessExif(Photo{
		ID:               id,
		Format:           format,
		Path:             fpath,
		FilePath:         path.Join(dir, fpath),
		OriginalFilePath: path.Join(dir, "o", fpath),
		SourcePath:       imgPath,
		Hash:             hash.Sum(nil),
		Width:            img.Bounds().Dx(),
		Height:           img.Bounds().Dy(),
		PlaceholderURI:   template.URL(generatePlaceholderURI(bytes.NewReader(buf))),
	}, x)

	err = downloadPhoto(p.FilePath, bytes.NewReader(buf))
	if err != nil {
		return Photo{}, err
	}

	err = downloadPhoto(p.OriginalFilePath, bytes.NewReader(buf))
	return p, err
}

func ProcessExif(img Photo, x *exif.Exif) Photo {
	mkTag, _ := x.Get(exif.Make)
	if mkTag != nil {
		img.MakeModel, _ = mkTag.StringVal()
	}

	modelTag, _ := x.Get(exif.Model)
	if modelTag != nil {
		model, _ := modelTag.StringVal()
		img.MakeModel += " " + model
	}

	img.DateTime, _ = x.DateTime()

	ssTag, _ := x.Get(exif.ShutterSpeedValue)
	if ssTag != nil {
		ssApexRat, _ := ssTag.Rat(0)
		ssApex, _ := ssApexRat.Float64()
		ss := math.Pow(2, -ssApex)

		if ss > 1 {
			img.ShutterSpeed = fmt.Sprintf("%.1fs", ss)
		} else {
			img.ShutterSpeed = fmt.Sprintf("1/%ds", int(math.Round(1/ss)))
		}
	}

	fTag, _ := x.Get(exif.FNumber)
	if fTag != nil {
		fRat, _ := fTag.Rat(0)
		img.FNumber = "ƒ/"
		if fRat.IsInt() {
			img.FNumber += fRat.RatString()
		} else {
			f, _ := fRat.Float64()
			img.FNumber += fmt.Sprintf("%.1f", f)
		}
	}

	isoTag, _ := x.Get(exif.ISOSpeedRatings)
	if isoTag != nil {
		iso, _ := isoTag.Int64(0)
		img.ISO = "ISO " + fmt.Sprint(iso)
	}

	return img
}

func base58Encode(src []byte) string {
	alphabet := []rune("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

	bytes := make([]byte, 0)

	leadingZeros := 0
	for _, b := range src {
		if b == 0 {
			leadingZeros++
		} else {
			break
		}
	}

	for _, b := range src {
		carry := int(b)
		for j := 0; carry != 0 || j < len(bytes); j++ {
			if j == len(bytes) {
				carry += 0
			} else {
				carry += int(bytes[j]) << 8
			}

			if j == len(bytes) {
				bytes = append(bytes, byte(carry%58))
			} else {
				bytes[j] = byte(carry % 58)
			}

			carry /= 58
		}
	}

	str := ""
	for i := 0; i < leadingZeros+len(bytes); i++ {
		if i < leadingZeros {
			str += string(alphabet[0])
		} else {
			str += string(alphabet[int(bytes[len(bytes)+leadingZeros-i-1])])
		}
	}

	return str
}
