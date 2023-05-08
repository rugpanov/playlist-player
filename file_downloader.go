package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func(out *os.File) {
		if err := out.Close(); err != nil {
			log.Fatal(err)
		}
	}(out)

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
