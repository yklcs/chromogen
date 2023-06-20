package photos

import "path"

type Photos struct {
	photos       []Photo
	OriginalsDir string
	Dir          string
	BucketURL    string
}

func NewPhotos(dir string) Photos {
	p := Photos{
		Dir:          dir,
		OriginalsDir: path.Join(dir, "o"),
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

func (ps *Photos) Append(p Photo) {
	ps.photos = append(ps.photos, p)
}

func (ps *Photos) Find(id string) *Photo {
	for _, p := range ps.photos {
		if p.ID == id {
			return &p
		}
	}

	return nil
}
