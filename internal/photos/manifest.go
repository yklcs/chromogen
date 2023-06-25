package photos

import (
	"encoding/json"
	"os"
	"sync"
)

type Manifest struct {
	Hashes     map[string]string `json:"hashes"`
	Read       []string          `json:"read"`
	Compressed []string          `json:"compressed"`
	lock       sync.Mutex
	path       string
}

func NewManifest(path string) *Manifest {
	m := Manifest{
		Hashes: make(map[string]string),
		lock:   sync.Mutex{},
		path:   path,
	}

	m.Save()

	return &m
}

func (m *Manifest) Run(fn func(*Manifest) error) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.loadWithoutLock()
	err := fn(m)
	m.saveWithoutLock()

	return err
}

func (m *Manifest) saveWithoutLock() error {
	f, err := os.Create(m.path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

func (m *Manifest) Save() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.saveWithoutLock()
}

func (m *Manifest) Load() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.loadWithoutLock()
}

func (m *Manifest) loadWithoutLock() error {
	data, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, m)
	return err
}

func (m *Manifest) IsRead(id string) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := m.loadWithoutLock()
	if err != nil {
		return false, err
	}

	for _, v := range m.Read {
		if v == id {
			return true, nil
		}
	}

	return false, nil
}

func (m *Manifest) IsCompressed(id string) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := m.loadWithoutLock()
	if err != nil {
		return false, err
	}

	for _, v := range m.Compressed {
		if v == id {
			return true, nil
		}
	}

	return false, nil
}

func (m *Manifest) FindId(filename string) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := m.loadWithoutLock()
	if err != nil {
		return false, err
	}

	for k := range m.Hashes {
		if k == filename {
			return true, nil
		}
	}

	return false, nil
}
