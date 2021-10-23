package data

import (
	"database/sql"
	"encoding/json"
)

type IllustRepository struct {
	db *sql.DB
}

func NewIllustRepository(db *sql.DB) *IllustRepository {
	return &IllustRepository{db: db}
}

func (repo *IllustRepository) Prepare() error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(`create table if not exists illusts(
		body text,
		id text generated always as (json_extract(body, '$.id')) virtual not null,
		r18 boolean generated always as (json_extract(body, '$.r18')) virtual not null
	)`)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = repo.db.Exec(`create index if not exists xid on illusts(id)`)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = repo.db.Exec(`create index if not exists xr18 on illusts(r18)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repo *IllustRepository) IsExists(illustId string) bool {
	row := repo.db.QueryRow(`select exists (select 1 from illusts where id = ? limit 1)`, illustId)
	var exists bool
	err := row.Scan(&exists)
	return err == nil && exists
}

func (repo *IllustRepository) Save(illust *Illust) error {
	bytes, err := json.Marshal(illust)
	if err != nil {
		return err
	}
	body := string(bytes)
	_, err = repo.db.Exec(`insert or replace into illusts values (?)`, body)
	return err
}

func (repo *IllustRepository) GetByID(illustId string) (*Illust, error) {
	row := repo.db.QueryRow(`select body from illusts where id = ?`, illustId)
	var body string
	err := row.Scan(&body)
	if err != nil {
		return nil, err
	}
	var illust Illust
	err = json.Unmarshal([]byte(body), &illust)
	if err != nil {
		return nil, err
	}
	return &illust, nil
}

func (repo *IllustRepository) GetRandom(r18 int, limit int) ([]Illust, error) {
	rows, err := repo.db.Query(`select body from illusts where r18 = ? order by random() limit ?`, r18, limit)
	if err != nil {
		return nil, err
	}
	return scanRows(rows)
}

func (repo *IllustRepository) Search(r18 int, q string, limit int) ([]Illust, error) {
	rows, err := repo.db.Query(
		"select body from illusts where r18 = ? and body like '%"+q+"%' order by random() limit ?", r18, limit)
	if err != nil {
		return nil, err
	}
	return scanRows(rows)
}

func (repo *IllustRepository) Delete(id string) error {
	_, err := repo.db.Exec(`delete from illusts where id = ?`, id)
	return err
}

func scanRows(rows *sql.Rows) ([]Illust, error) {
	illusts := make([]Illust, 0)
	for rows.Next() {
		var body string
		err := rows.Scan(&body)
		if err != nil {
			return nil, err
		}
		var illust Illust
		err = json.Unmarshal([]byte(body), &illust)
		if err != nil {
			return nil, err
		}
		illusts = append(illusts, illust)
	}
	return illusts, nil
}
