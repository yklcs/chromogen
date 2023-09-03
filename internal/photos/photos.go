package photos

import (
	"database/sql"
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/yklcs/panchro/internal/photo"
	_ "modernc.org/sqlite"
)

type Photos struct {
	DB *sql.DB
}

func (ps *Photos) Load(dbpath string) (bool, error) {
	if _, err := os.Stat(dbpath); errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	db, err := sql.Open("sqlite", dbpath)
	if err != nil {
		return false, err
	}
	ps.DB = db
	return true, nil
}

func (ps *Photos) Init() error {
	_, err := ps.DB.Exec(`
		CREATE TABLE IF NOT EXISTS photos(
			id TEXT PRIMARY KEY,
			url TEXT,
			path TEXT,
			source_path TEXT,
			format TEXT,
			hash BLOB,
			placeholder_uri TEXT,
			width INTEGER,
			height INTEGER,
			
			exif_datetime DATETIME,
			exif_makemodel TEXT,
			exif_shutterspeed TEXT,
			exif_fnumber TEXT,
			exif_iso TEXT,
			exif_lensmakemodel TEXT,
			exif_focallength TEXT,
			exif_subjectdistance TEXT
		);
	`)

	return err
}

func (ps Photos) Add(p photo.Photo) {
	_, err := ps.DB.Exec(`
		INSERT OR REPLACE INTO photos
		(
			id,
			url,
			path,
			source_path,
			format,
			hash,
			placeholder_uri,
			width, 
			height,
			exif_datetime,
			exif_makemodel,
			exif_shutterspeed,
			exif_fnumber,
			exif_iso,
			exif_lensmakemodel,
			exif_focallength,
			exif_subjectdistance
		)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17);
		`,
		p.ID,
		p.URL,
		p.Path,
		p.SourcePath,
		p.Format,
		p.Hash,
		p.PlaceholderURI,
		p.Width,
		p.Height,
		p.Exif.DateTime,
		p.Exif.MakeModel,
		p.Exif.ShutterSpeed,
		p.Exif.FNumber,
		p.Exif.ISO,
		p.Exif.LensMakeModel,
		p.Exif.FocalLength,
		p.Exif.SubjectDistance,
	)
	if err != nil {
		log.Println(err)
	}

}

func (ps Photos) IDs() []string {
	var ids []string

	rows, _ := ps.DB.Query(
		`SELECT id FROM photos ORDER BY ROWID DESC;`,
	)

	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids
}

func (ps Photos) Len() int {
	return len(ps.IDs())
}

func (ps Photos) Get(id string) (photo.Photo, error) {
	var p photo.Photo
	p.Exif = &photo.Exif{}

	row := ps.DB.QueryRow(`
	SELECT * FROM photos
	WHERE id = ?;`, id)
	err := row.Scan(
		&p.ID,
		&p.URL,
		&p.Path,
		&p.SourcePath,
		&p.Format,
		&p.Hash,
		&p.PlaceholderURI,
		&p.Width,
		&p.Height,
		&p.Exif.DateTime,
		&p.Exif.MakeModel,
		&p.Exif.ShutterSpeed,
		&p.Exif.FNumber,
		&p.Exif.ISO,
		&p.Exif.LensMakeModel,
		&p.Exif.FocalLength,
		&p.Exif.SubjectDistance,
	)

	return p, err
}

func (ps Photos) Delete(id string) error {
	_, err := ps.DB.Exec(
		`DELETE FROM photos
		WHERE id = ?
		`, id,
	)

	return err
}
