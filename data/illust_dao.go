package data

import (
	"database/sql"
	"strings"
)

const table_illusts = `
CREATE TABLE IF NOT EXISTS illusts (
    id          TEXT,
    title       TEXT,
    description TEXT,
    author_id   TEXT,
    author      TEXT,
    account     TEXT,
    r18         BOOLEAN,
    tags        TEXT,
    create_date DATETIME,
    upload_date DATETIME,
    PRIMARY KEY (id)
)
`

const table_illust_tags = `
CREATE TABLE IF NOT EXISTS illust_tags (
    illust_id   TEXT,
    tag         TEXT,
    PRIMARY KEY (illust_id, tag)
)
`

const table_illust_pages = `
CREATE TABLE IF NOT EXISTS illust_pages (
    illust_id   TEXT,
    p           INTEGER,
    width		INTEGER,
    height      INTEGER,
    thumb       TEXT,
    small       TEXT,
    regular     TEXT,
    original    TEXT,
    PRIMARY KEY (illust_id, p)
)
`

type IllustDao struct {
	db *sql.DB
}

func NewIllustDao(db *sql.DB) (*IllustDao, error) {
	err := createTables(db)
	if err != nil {
		return nil, err
	}
	return &IllustDao{db: db}, nil
}

func createTables(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(table_illusts)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(table_illust_tags)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(table_illust_pages)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (dao *IllustDao) CheckExists(id string) bool {
	rows, err := dao.db.Query("SELECT 1 FROM illusts WHERE id = ?", id)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func (dao *IllustDao) FindByID(id string) (*Illust, error) {
	query := `SELECT 
        illusts.id, illusts.title, illusts.description, illusts.author_id, illusts.author,
        illusts.account, illusts.r18,illusts.create_date,illusts.upload_date,
        tags.tag,
        pages.p, pages.width, pages.height, pages.thumb, pages.small, pages.regular, pages.original
    FROM illusts 
        LEFT JOIN illust_tags AS tags ON tags.illust_id = illusts.id
        LEFT JOIN illust_pages AS pages ON pages.illust_id = illusts.id
    WHERE illusts.id = ?
    `
	rows, err := dao.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	illust := Illust{}
	for rows.Next() {
		var (
			tag      string
			p        int
			width    int
			height   int
			thumb    string
			small    string
			regular  string
			original string
		)
		rows.Scan(&illust.ID, &illust.Title, &illust.Description, &illust.AuthorID, &illust.Account,
			&illust.R18, &illust.CreateDate, &illust.UploadDate, &tag, &p, &width, &height, &thumb,
			&small, &regular, &original)
	}
	return nil, nil
}

func (dao *IllustDao) Save(illust *Illust) error {
	tx, err := dao.db.Begin()
	if err != nil {
		return err
	}

	query := `INSERT OR REPLACE INTO illusts (
        id, title, description, author_id, author, account, r18, tags, create_date, upload_date
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(query,
		illust.ID,
		illust.Title,
		illust.Description,
		illust.AuthorID,
		illust.Author,
		illust.Account,
		illust.R18,
		strings.Join(illust.Tags, ","),
		illust.CreateDate,
		illust.UploadDate)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = "INSERT OR REPLACE INTO illust_tags (illust_id, tag) VALUES "
	values := make([]string, 0, len(illust.Tags))
	args := make([]interface{}, 0, len(illust.Tags)*2)
	for _, tag := range illust.Tags {
		values = append(values, "(?, ?)")
		args = append(args, illust.ID, tag)
	}
	query = query + strings.Join(values, ",")
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = "INSERT OR REPLACE INTO illust_pages (illust_id, p, width, height, thumb, small, regular, original) VALUES"
	values = make([]string, 0, len(illust.Pages))
	args = make([]interface{}, 0, len(illust.Pages)*8)
	for p, page := range illust.Pages {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args, illust.ID, p, page.Width, page.Height, page.Thumb, page.Small, page.Regular, page.Original)
	}
	query = query + strings.Join(values, ",")
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (dao *IllustDao) Random(r18 int, limit int) ([]Illust, error) {
	panic("unimplement yet")
}

func (dao *IllustDao) RandomSearch(r18 int, q string, limit int) ([]Illust, error) {
	panic("unimplement yet")
}
