package photos

import (
	"encoding/json"
	"errors"

	"github.com/yklcs/panchro/internal/photo"
	bolt "go.etcd.io/bbolt"
)

type Photos struct {
	DB *bolt.DB
}

func (ps *Photos) MarshalJSON() ([]byte, error) {
	var pslice []photo.Photo
	for _, id := range ps.IDs() {
		p, err := ps.Get(id)
		if err != nil {
			return nil, err
		}
		pslice = append(pslice, p)
	}
	return json.Marshal(pslice)
}

var IndexKey = []byte("index")
var PhotosKey = []byte("photos")

func (ps *Photos) Init() error {
	err := ps.DB.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(IndexKey)
		tx.CreateBucketIfNotExists(PhotosKey)
		return nil
	})
	return err
}

func (ps Photos) Add(p photo.Photo) {
	ps.DB.Update(func(tx *bolt.Tx) error {
		// Add p.ID to index
		indexBucket := tx.Bucket(IndexKey)
		id, _ := indexBucket.NextSequence()
		indexBucket.Put(uint64ToByteSlice(id), []byte(p.ID))

		// Add p
		pb, err := json.Marshal(p)
		if err != nil {
			return err
		}
		photosBucket := tx.Bucket(PhotosKey)
		photosBucket.Put([]byte(p.ID), pb)

		return nil
	})
}

func (ps Photos) Set(p photo.Photo) {
	ps.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(PhotosKey)
		pb, err := json.Marshal(p)
		if err != nil {
			return err
		}
		b.Put([]byte(p.ID), pb)
		return nil
	})
}

func (ps Photos) IDs() []string {
	ids := []string{}

	ps.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(IndexKey)
		b.ForEach(func(k, v []byte) error {
			ids = append(ids, string(v))
			return nil
		})
		return nil
	})

	return ids
}

func (ps Photos) Len() int {
	return len(ps.IDs())
}

func (ps Photos) Get(id string) (photo.Photo, error) {
	var p photo.Photo
	err := ps.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(PhotosKey)
		pb := b.Get([]byte(id))
		if pb == nil {
			return errors.New("photo does not exist")
		}
		err := json.Unmarshal(pb, &p)
		return err
	})

	return p, err
}

func (ps Photos) Delete(id string) error {
	err := ps.DB.Update(func(tx *bolt.Tx) error {
		// Delete from photos
		photosBucket := tx.Bucket(PhotosKey)
		err := photosBucket.Delete([]byte(id))
		if err != nil {
			return err
		}

		// Delete from index
		indexBucket := tx.Bucket(IndexKey)
		c := indexBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if id == string(v) {
				indexBucket.Delete(k)
			}
		}

		return nil
	})

	return err
}

func uint64ToByteSlice(u uint64) []byte {
	return []byte{
		byte(0xff & u),
		byte(0xff & (u >> 8)),
		byte(0xff & (u >> 16)),
		byte(0xff & (u >> 24)),
		byte(0xff & (u >> 32)),
		byte(0xff & (u >> 40)),
		byte(0xff & (u >> 48)),
		byte(0xff & (u >> 56)),
	}
}
