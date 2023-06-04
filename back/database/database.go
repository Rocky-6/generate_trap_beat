package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Chord struct {
	Degree_name string
	Function    string
}

func DispChord() [4]string {
	db, err := sql.Open("sqlite3", "./trap.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var cp [4]string
	var degree_name string
	var function string

	sql := "SELECT * FROM chord WHERE function='T' ORDER BY RANDOM() LIMIT 1;"

	err = db.QueryRow(sql).Scan(&degree_name, &function)
	if err != nil {
		log.Fatal(err)
	}

	cp[0] = degree_name

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

		err = db.QueryRow(sql).Scan(&degree_name, &function)
		if err != nil {
			log.Fatal(err)
		}

		cp[i] = degree_name
	}

	return cp
}
