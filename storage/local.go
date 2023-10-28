package storage

import (
	"io"
	"net/http"
	"os"
	"path"
)

type LocalStorage struct {
	dir    string
	prefix string
}

func NewLocalStorage(dir, prefix string) (*LocalStorage, error) {
	err := os.MkdirAll(path.Join(dir, prefix), 0755)
	if err != nil {
		return nil, err
	}

	return &LocalStorage{dir: dir, prefix: prefix}, nil
}

func (s *LocalStorage) Upload(r io.Reader, fpath string) (string, error) {
	fpathjoined := path.Join(s.dir, s.prefix, fpath)
	f, err := os.Create(fpathjoined)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.ReadFrom(r)
	if err != nil {
		return "", err
	}

	return path.Join(s.prefix, fpath), err
}

func (s *LocalStorage) Delete(fpath string) error {
	fpathjoined := path.Join(s.dir, fpath)
	err := os.Remove(fpathjoined)
	return err
}

func (s *LocalStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join(s.dir, r.URL.Path)
	http.ServeFile(w, r, filepath)
}

func (s LocalStorage) Backend() string {
	return "local"
}
