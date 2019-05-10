package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
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
VALUES (2, '%s', '%s', '%s', %d, '%s' )`, path, album, name, date, note)
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

func GetThumbnailsHandler() []ImageResponse {
	var albums []string
	db := getDB()
	defer db.Close()

	sqlAlbumStatement := `
SELECT DISTINCT album FROM Images;`
	row, err := db.Query(sqlAlbumStatement)
	defer row.Close()
	for row.Next() {
		var albumName string
		err = row.Scan(&albumName)
		if err != nil {
			panic(err)
		}
		albums = append(albums, albumName)
	}

	var response []ImageResponse
	for i := 0; i < len(albums); i++ {
		sqlStatement := fmt.Sprintf(`SELECT uuid, path, album, name, date FROM Images WHERE album='%s' LIMIT 1`, albums[i])

		var uid int
		var path string
		var album string
		var name string
		var date int

		row := db.QueryRow(sqlStatement)
		err = row.Scan(&uid, &path, &album, &name, &date)
		response = append(response, ImageResponse{uid, path, album, name, date})
	}

	return response
}

func GetAlbumNamesHandler() []string {
	var albums []string
	db := getDB()
	defer db.Close()

	sqlAlbumStatement := `
SELECT DISTINCT album FROM Images;`
	row, err := db.Query(sqlAlbumStatement)
	defer row.Close()
	for row.Next() {
		var albumName string
		err = row.Scan(&albumName)
		if err != nil {
			panic(err)
		}
		albums = append(albums, albumName)
	}

	return albums
}

func GetAlbumHandler(albumName string) []ImageResponse {

	db := getDB()
	defer db.Close()
	sqlStatement := fmt.Sprintf(`SELECT uuid, path, album, name, date FROM Images WHERE album='%s';`, strings.Title(albumName))
	row2, err := db.Query(sqlStatement)
	defer row2.Close()

	var response []ImageResponse
	for row2.Next() {
		var uid int
		var path string
		var album string
		var name string
		var date int

		err = row2.Scan(&uid, &path, &album, &name, &date)
		if err != nil {
			panic(err)
		}
		response = append(response, ImageResponse{uid, path, album, name, date})
	}
	return response

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
