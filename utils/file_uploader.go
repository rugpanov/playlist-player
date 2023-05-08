package utils

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// UploadFile :: taken from https://tutorialedge.net/golang/go-file-upload-tutorial/
func UploadFile(w http.ResponseWriter, r *http.Request) (string, error) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader, so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer func(file multipart.File) {
		if err := file.Close(); err != nil {
			log.Fatal(err)
			return
		}
	}(file)
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	newUuid := uuid.New()
	songFile, err := os.Create("music/uploaded-" + newUuid.String() + ".mp3")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer func(songFile *os.File) {
		if err := songFile.Close(); err != nil {
			log.Fatal(err)
			return
		}
	}(songFile)

	// read all the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// write this byte array to our temporary file
	_, err = songFile.Write(fileBytes)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	// return that we have successfully uploaded our file!
	_, err = fmt.Fprintf(w, "Successfully Uploaded File\n")
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return songFile.Name(), nil
}
