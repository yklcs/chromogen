package config

import (
	"encoding/json"
	"os"
)

type viewMode string

const (
	GalleryMode viewMode = "gallery"
	GridMode    viewMode = "grid"
)

type Config struct {
	Title           string                 `json:"title"`
	Root            string                 `json:"root"`
	DefaultViewMode viewMode               `json:"default_view_mode"`
	Theme           string                 `json:"theme"`
	ThemeConfig     map[string]interface{} `json:"theme_config"`
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
