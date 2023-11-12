package service

import (
	"context"
	"database/sql"

	"github.com/Rocky-6/trap/model"
	"github.com/Rocky-6/trap/repository"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite3Client struct {
	client *sql.DB
}

func NewSqliteClient(src string) (repository.DBRepository, error) {
	db, err := sql.Open("sqlite3", src)
	if err != nil {
		return nil, err
	}
	return &sqlite3Client{client: db}, err
}

func (c *sqlite3Client) Scan(ctx context.Context) ([]model.ChordInfomation, error) {
	var (
		degreeName string
		function   string
	)
	chordInfomation := make([]model.ChordInfomation, 0, 4)

	sql := "SELECT * FROM chord WHERE function='T' ORDER BY RANDOM() LIMIT 1;"

	if err := c.client.QueryRow(sql).Scan(&degreeName, &function); err != nil {
		return nil, err
	}
	chordInfomation = append(chordInfomation, model.ChordInfomation{
		DegreeName: degreeName,
		Function:   function,
	})

	for i := 1; i < 4; i++ {
		switch function {
		case "T":
			sql = "SELECT * FROM chord WHERE (function='D' OR function='S' OR function='SM') ORDER BY RANDOM() LIMIT 1;"
		case "D":
			sql = "SELECT * FROM chord WHERE function='T' ORDER BY RANDOM() LIMIT 1;"
		case "S":
			sql = "SELECT * FROM chord WHERE (function='T' OR function='D' OR function='SM') ORDER BY RANDOM() LIMIT 1;"
		case "SM":
			sql = "SELECT * FROM chord WHERE (function='T' OR function='D') ORDER BY RANDOM() LIMIT 1;"
		}

		if err := c.client.QueryRow(sql).Scan(&degreeName, &function); err != nil {
			return nil, err
		}
		chordInfomation = append(chordInfomation, model.ChordInfomation{
			DegreeName: degreeName,
			Function:   function,
		})
	}

	return chordInfomation, nil
}
