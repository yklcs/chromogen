package config

import (
	"encoding/json"
	"os"
)

type ViewMode string

const (
	GalleryMode ViewMode = "gallery"
	GridMode    ViewMode = "grid"
)

type Config struct {
	Title           string   `json:"title"`
	Root            string   `json:"root"`
	DefaultViewMode ViewMode `json:"default_view_mode"`
	StaticDir       string
}

func ReadConfig(path string) (*Config, error) {
	c := &Config{}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}

	c.StaticDir = "static"
	return c, nil
}
