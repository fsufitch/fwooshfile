package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	apiurlPtr := flag.String("apiurl", "http://localhost:8080",
		"the URL at which the filebounce API is located")
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

	fmt.Printf("Uploading '%s' to '%s'...\n", filepath, apiurl)
}
