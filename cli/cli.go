package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func checkAPIAlive(apiurl string) error {
	if apiurl[len(apiurl)-1] != "/"[0] {
		apiurl += "/"
	}

	statusurl := apiurl + "api/status"
	resp, err := http.Get(statusurl)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Non-OK status code: %d", resp.StatusCode)
	}

	statusBody := [2]byte{}

	_, _ = resp.Body.Read(statusBody[:])
	if !bytes.Equal(statusBody[:], []byte("OK")) {
		return fmt.Errorf("Non-OK query body: %s", string(statusBody[:]))
	}

	return nil
}

func main() {
	apiurlPtr := flag.String("apiurl", "http://localhost:8080",
		"the URL at which the filebounce API is located")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-h] [-apiurl APIURL] FILE\n", os.Args[0])
	}
	flag.Parse()

	apiurl := *apiurlPtr
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please specify a file to upload")
		return
	}

	filepath := args[0]
	fileinfo, err := os.Stat(filepath)
	if err != nil {
		fmt.Printf("File %s does not exist\n", filepath)
		return
	}
	if fileinfo.IsDir() {
		fmt.Printf("Cannot upload file %s; it is a directory\n", filepath)
		return
	}

	if err = checkAPIAlive(apiurl); err != nil {
		fmt.Printf("API not alive: %s\n", err.Error())
		return
	}

	fmt.Printf("Uploading '%s' to '%s'...\n", filepath, apiurl)
}
