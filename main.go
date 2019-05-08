package main

import (
	"Projects/Gati/controllers"

	"github.com/gorilla/mux"

	"fmt"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/album/{album}", controllers.GetAlbum).Methods("GET")
	router.HandleFunc("/api/image/{id}", controllers.GetImage).Methods("GET")
	router.HandleFunc("/api/upload", controllers.UploadImage).Methods("POST")

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
