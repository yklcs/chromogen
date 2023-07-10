package render

import (
	"io"
	"io/fs"
	"os"
	"path"
)

func CopyFS(fsys fs.FS, dst string) error {
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	return fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			err = os.MkdirAll(path.Join(dst, p), 0755)
			if err != nil {
				return err
			}
		} else {
			fsrc, err := fsys.Open(p)
			if err != nil {
				return err
			}
			defer fsrc.Close()

			fdst, err := os.Create(path.Join(dst, p))
			if err != nil {
				return err
			}
			defer fdst.Close()

			io.Copy(fdst, fsrc)
		}

		return nil
	})
}
