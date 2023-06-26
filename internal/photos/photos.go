package photos

import (
	"github.com/yklcs/panchro/internal/photo"
	"github.com/yklcs/panchro/storage"
	"golang.org/x/exp/slices"
)

type Photos struct {
	photos   []photo.Photo
	manifest *Manifest
	store    storage.Storage
}

func NewPhotos() Photos {
	p := Photos{
		manifest: NewManifest("manifest.json"),
	}

	return p
}

func (ps Photos) Slice() []photo.Photo {
	return ps.photos
}

func (ps Photos) Len() int {
	return len(ps.photos)
}

func (ps *Photos) Get(index int) *photo.Photo {
	return &ps.photos[index]
}

func (ps *Photos) Add(p photo.Photo) {
	ps.manifest.Run(func(m *Manifest) error {
		m.Hashes[p.SourcePath] = p.ID
		m.Read = append(ps.manifest.Read, p.ID)
		return nil
	})
	pos, _ := slices.BinarySearchFunc(
		ps.photos, p, func(p1, p2 photo.Photo) int {
			return -p1.Exif.DateTime.Compare(p2.Exif.DateTime)
		},
	)
	ps.photos = slices.Insert(ps.photos, pos, p)
}

func (ps *Photos) Find(id string) *photo.Photo {
	for _, p := range ps.photos {
		if p.ID == id {
			return &p
		}
	}

	return nil
}
