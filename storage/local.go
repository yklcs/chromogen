package storage

import (
	"io"
	"net/http"
	"os"
	"path"
)

type LocalStorage struct {
	dir string
}

func NewLocalStorage(dir string) (*LocalStorage, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	return &LocalStorage{dir: dir}, nil
}

func (s *LocalStorage) Upload(r io.Reader, fpath string) (string, error) {
	fpathjoined := path.Join(s.dir, fpath)
	f, err := os.Create(fpathjoined)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.ReadFrom(r)
	if err != nil {
		return "", err
	}

	return fpath, err
}

func (s *LocalStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join(s.dir, r.URL.Path)
	http.ServeFile(w, r, filepath)
}

func (s LocalStorage) Backend() string {
	return "local"
}

type localReader struct {
	dir   string
	fpath string
}

func (r *localReader) Read(b []byte) (int, error) {
	return 0, nil
}

func newLocalReader(s *LocalStorage, fpath string) localReader {
	return localReader{
		dir:   s.dir,
		fpath: fpath,
	}
}
