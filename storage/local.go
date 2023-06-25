package storage

import "net/http"

type LocalStorage struct {
}

func (s *LocalStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
