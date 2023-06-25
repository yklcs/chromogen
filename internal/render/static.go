package render

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path"

	"gocloud.dev/blob"
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

			fdst, err := os.Create(path.Join(dst, p))
			if err != nil {
				return err
			}

			io.Copy(fdst, fsrc)
		}

		return nil
	})
}

func WriteFS(fsys fs.FS, bucket *blob.Bucket) error {
	return fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// err = os.MkdirAll(path.Join(dst, p), 0755)
			// if err != nil {
			// return err
			// }
		} else {
			fsrc, err := fsys.Open(p)
			if err != nil {
				return err
			}

			w, err := bucket.NewWriter(context.Background(), p, nil)
			if err != nil {
				return err
			}
			defer w.Close()
			// fdst, err := os.Create(path.Join(dst, p))
			// if err != nil {
			// return err
			// }

			io.Copy(w, fsrc)
		}

		return nil
	})
}
