package photo

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math"

	"github.com/rwcarlsen/goexif/exif"
)

func processExif(x *exif.Exif) *Exif {
	ex := Exif{}

	mkTag, _ := x.Get(exif.Make)
	if mkTag != nil {
		ex.MakeModel, _ = mkTag.StringVal()
	}

	oTag, _ := x.Get(exif.Orientation)
	fmt.Println(oTag)

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
