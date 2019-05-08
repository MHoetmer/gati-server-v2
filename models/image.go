package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	pq "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Postgres"
	dbname   = "gatiway"
)

type ImageResponse struct {
	Uid   int
	Path  string
	Album string
	Name  string
	Date  int
}

func SaveFilePath(path, album, name, note string) {
	db := getDB()
	defer db.Close()
	date := time.Now().Unix()

	sqlStatement := fmt.Sprintf(`
INSERT INTO Images (uuid, path, album, name, date, note)
VALUES (31, '%s', '%s', '%s', %d, '%s' )`, path, album, name, date, note)
	fmt.Println("sqlStatement", sqlStatement)
	_, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
	fmt.Println("insert result: ", sqlStatement)
}

func GetImageByUuid(uuid int) *ImageResponse {
	var uid int
	var path string
	var album string
	var name string
	var date int

	db := getDB()
	defer db.Close()

	sqlStatement := `
SELECT uuid, path, album, name, date FROM Images WHERE uuid=$1;`
	row := db.QueryRow(sqlStatement, uuid)
	//fmt.Println(row)
	switch err := row.Scan(&uid, &path, &album, &name, &date); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
	case nil:
		return &ImageResponse{uid, path, album, name, date}
	default:
		panic(err)
	}

	return &ImageResponse{}
}

func GetAlbum(albumName string) (albumResponse []*ImageResponse) {

	db := getDB()
	defer db.Close()

	sqlStatement := `
SELECT uuid, path, album, name, date FROM Images WHERE album=$1;`
	row := db.QueryRow(sqlStatement, albumName)
	//fmt.Println(row)

	if err := row.Scan(pq.Array(&albumResponse)); err != nil {
		log.Fatal(err)
	}
	return
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
