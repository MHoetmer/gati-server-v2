package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Postgres"
	dbname   = "gatiway"
)

type Response struct {
	Uid   int
	Path  string
	Album string
	Name  string
	Date  int
	Note  string
}

func GetImage(uuid int) *Response {
	var uid int
	var path string
	var album string
	var name string
	var date int
	var note string

	db := getDB()
	defer db.Close()

	sqlStatement := `
SELECT * FROM Images WHERE uuid=$1;`
	row := db.QueryRow(sqlStatement, uuid)
	//fmt.Println(row)
	switch err := row.Scan(&uid, &path, &album, &name, &date, &note); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
	case nil:
		return &Response{uid, path, album, name, date, note}
	default:
		panic(err)
	}

	return &Response{}
}

func getDB() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return db
}
