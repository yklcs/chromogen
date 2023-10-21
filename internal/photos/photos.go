package photos

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type Photos struct {
	DB *sql.DB
}

func (ps *Photos) Init() error {
	_, err := ps.DB.Exec(`
		CREATE TABLE IF NOT EXISTS photos(
			id TEXT PRIMARY KEY,
			url TEXT NOT NULL,
			path TEXT NOT NULL,
			format TEXT NOT NULL,
			hash BLOB NOT NULL,
			placeholder_uri TEXT NOT NULL,
			width INTEGER NOT NULL,
			height INTEGER NOT NULL,
			exif_datetime DATETIME NOT NULL,
			exif_makemodel TEXT NOT NULL,
			exif_shutterspeed TEXT NOT NULL,
			exif_fnumber TEXT NOT NULL,
			exif_iso TEXT NOT NULL,
			exif_lensmakemodel TEXT NOT NULL,
			exif_focallength TEXT NOT NULL,
			exif_subjectdistance TEXT NOT NULL
		);
	`)

	return err
}

func (ps Photos) Set(p *Photo) {
	if p == nil {
		return
	}

	if p.Exif == nil {
		p.Exif = &Exif{}
	}

	_, err := ps.DB.Exec(`
		INSERT OR REPLACE INTO photos
		(
			id,
			url,
			path,
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
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);
		`,
		p.ID,
		p.URL,
		p.Path,
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

func (ps Photos) Get(id string) (*Photo, error) {
	var p Photo
	p.Exif = &Exif{}

	row := ps.DB.QueryRow(`
	SELECT * FROM photos
	WHERE id = ?;`, id)
	err := row.Scan(
		&p.ID,
		&p.URL,
		&p.Path,
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

	return &p, err
}

func (ps Photos) Delete(id string) error {
	_, err := ps.DB.Exec(`
		DELETE FROM photos
		WHERE id = ?`, id,
	)

	return err
}
