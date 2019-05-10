package controllers

import (
	"Projects/Gati/models"
	u "Projects/Gati/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	directory := "/Users/mel/Projects/javascript/Gati/src/photos"
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	album := r.FormValue("Album")
	albumDirectory := directory + "/" + album
	name := r.FormValue("Name")
	note := r.FormValue("Note")

	file, handler, err := r.FormFile("File")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	fmt.Println("album", album)
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	strippedName := strings.Split(handler.Filename, ".")
	pngName := fmt.Sprintf("%s.png", strippedName[0])
	fmt.Println("New name:", pngName)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(albumDirectory, "*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")

	oldName := fmt.Sprintf("%s", tempFile.Name())
	newName := albumDirectory + "/" + pngName

	err = os.Rename(oldName, newName)
	if err != nil {
		log.Fatal(err)
	}

	relativePath := strings.Split(newName, "src")
	fmt.Println("path", relativePath[1])

	models.SaveFilePath(relativePath[1], album, name, note)
}
func GetThumbnails(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	data := models.GetThumbnailsHandler()
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
func GetImage(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	fmt.Println("params", params)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		//The passed path parameter is not an integer
		u.Respond(w, u.Message(false, "There was an error in your request"))
		return
	}

	data := models.GetImageByUuid(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

func GetAlbum(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	album := params["album"]

	data := models.GetAlbumHandler(album)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

func GetAlbumNames(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	data := models.GetAlbumNamesHandler()
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
