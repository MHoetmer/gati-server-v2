package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Postgres"
	dbname   = "gatiway"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/image", getImageHandler).Methods("GET")
	fmt.Println("Serving at localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getImageHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Path[len("/image/"):]
	i, _ := strconv.Atoi(uid)
	//fmt.Println("image", i)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	resp := getImage(db, i)
	fmt.Println("resp", resp)

	f, err := os.Open("./photos/georgie/18722458_104334886830863_2699040101157044224_n.jpg")
	if err != nil {
		fmt.Println("could not open file", err)
	}
	// Read the entire JPG file into memory.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	//fmt.Println("content", content)

	// Set the Content Type header.
	w.Header().Set("Content-Type", "image/jpeg")
	//w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Write the image to the response.

	w.Write(content)

	fmt.Fprintf(w, "You are on the home page")
}

type Response struct {
	uid  int
	path string
	name string
	date int
	note string
}

func getImage(db *sql.DB, uuid int) Response {
	var uid int
	var path string
	var name string
	var date int
	var note string

	sqlStatement := `
SELECT * FROM Images WHERE uuid=$1;`
	row := db.QueryRow(sqlStatement, uuid)
	//fmt.Println(row)
	switch err := row.Scan(&uid, &path, &name, &date, &note); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
	case nil:
		return Response{uid, path, name, date, note}
	default:
		panic(err)
	}
	return Response{}
}

func saveImage(db *sql.DB) {
	sqlStatement := `
	INSERT INTO images (uuid, path, name, date, note)
	VALUES (10, './gati/photos/georgie/18722458_104334886830863_2699040101157044224_n.jpg', 'Bikers', 1557259719, 'Liever lui dan moe')`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}
