package photos

import (
	"path"

	"gocloud.dev/blob"
	"golang.org/x/exp/slices"
)

type Photos struct {
	photos       []Photo
	OriginalsDir string
	Dir          string
	manifest     *Manifest
	Bucket       *blob.Bucket
}

func NewPhotos(dir string, bucket *blob.Bucket) Photos {
	p := Photos{
		Dir:          dir,
		OriginalsDir: path.Join(dir, "o"),
		manifest:     NewManifest(path.Join(dir, "manifest.json")),
		Bucket:       bucket,
	}

	return p
}

func (ps Photos) Slice() []Photo {
	return ps.photos
}

func (ps Photos) Len() int {
	return len(ps.photos)
}

func (ps *Photos) Get(index int) *Photo {
	return &ps.photos[index]
}

func (ps *Photos) Add(p Photo) {
	ps.manifest.Run(func(m *Manifest) error {
		m.Hashes[p.SourcePath] = p.ID
		m.Read = append(ps.manifest.Read, p.ID)
		return nil
	})
	pos, _ := slices.BinarySearchFunc(
		ps.photos, p, func(p1, p2 Photo) int {
			return -p1.Exif.DateTime.Compare(p2.Exif.DateTime)
		},
	)
	ps.photos = slices.Insert(ps.photos, pos, p)
}

func (ps *Photos) Find(id string) *Photo {
	for _, p := range ps.photos {
		if p.ID == id {
			return &p
		}
	}

	return nil
}
