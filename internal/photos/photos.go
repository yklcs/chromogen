package photos

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
	"github.com/yklcs/panchro/internal/photo"
	"golang.org/x/exp/slices"
)

type Photos struct {
	db *badger.DB
}

var IndexKey = []byte{0}

func NewPhotos(db *badger.DB) Photos {
	db.Update(func(txn *badger.Txn) error {
		mapb, err := json.Marshal([]string{})
		if err != nil {
			return err
		}
		if _, err := txn.Get(IndexKey); err == badger.ErrKeyNotFound {
			txn.Set(IndexKey, mapb)
		}

		return nil
	})

	return Photos{
		db: db,
	}
}

func (ps Photos) Add(val photo.Photo) {
	var ids []string
	ps.db.Update(func(txn *badger.Txn) error {
		valb, err := json.Marshal(val)
		if err != nil {
			return err
		}
		txn.Set([]byte(val.ID), valb)

		idsitem, err := txn.Get(IndexKey)
		if err != nil {
			return err
		}
		err = idsitem.Value(func(val []byte) error {
			err := json.Unmarshal(val, &ids)
			return err
		})
		if err != nil {
			return err
		}

		idsb, err := json.Marshal(append([]string{val.ID}, ids...))
		if err != nil {
			return err
		}

		err = txn.Set(IndexKey, idsb)
		return err
	})
}

func (ps Photos) Set(val photo.Photo) {
	ps.db.Update(func(txn *badger.Txn) error {
		valb, err := json.Marshal(val)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(val.ID), valb)
		return err
	})
}

func (ps Photos) IDs() []string {
	ids := []string{}
	ps.db.View(func(txn *badger.Txn) error {
		idsitem, err := txn.Get(IndexKey)
		if err != nil {
			return err
		}
		err = idsitem.Value(func(val []byte) error {
			err := json.Unmarshal(val, &ids)
			return err
		})
		return err
	})

	return ids
}

func (ps Photos) Len() int {
	return len(ps.IDs())
}

func (ps Photos) Get(id string) (photo.Photo, error) {
	var p photo.Photo
	err := ps.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			err := json.Unmarshal(val, &p)
			return err
		})
		return err
	})
	return p, err
}

func (ps Photos) Delete(id string) (photo.Photo, error) {
	var p photo.Photo
	ids := []string{}
	err := ps.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(id))
		if err != nil {
			return err
		}

		idsitem, err := txn.Get(IndexKey)
		if err != nil {
			return err
		}
		err = idsitem.Value(func(val []byte) error {
			err := json.Unmarshal(val, &ids)
			return err
		})
		if err != nil {
			return err
		}

		i := slices.Index(ids, id)
		ids = slices.Delete(ids, i, i+1)

		idsb, err := json.Marshal(ids)
		if err != nil {
			return err
		}
		err = txn.Set(IndexKey, idsb)

		return err
	})
	return p, err
}
