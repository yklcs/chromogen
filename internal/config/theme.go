package config

import (
	"io/fs"
	"os"

	"github.com/yklcs/panchro/web"
)

func LoadTheme(conf *Config) fs.FS {
	var themeFS fs.FS
	if conf.Theme == "" {
		return web.Content
	} else {
		themeFS = os.DirFS(conf.Theme)
	}
	return themeFS
}
