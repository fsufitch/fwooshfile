package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsufitch/fwooshfile/client"
)

func main() {
	apiurlPtr := flag.String("apiurl", "http://localhost:8080/",
		"the URL at which the filebounce API is located")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-h] [-apiurl APIURL] FILE\n", os.Args[0])
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please specify a file to upload")
		return
	}
	filepath := args[0]

	transferClient := client.TransferClient{
		BaseURL: *apiurlPtr,
		Name:    "bouncecli",
	}

	fileTransfer := client.FileTransfer{
		Path: filepath,
	}

	transferID, err := fileTransfer.StartTransfer(transferClient)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing transfer: %s", err)
	}

	fmt.Printf("Upload registered. Download at: %s\n", transferClient.BaseURL+"d/"+transferID)

	fmt.Print("Push Enter to start upload: ")
	var input string
	fmt.Scanln(&input)
	fmt.Println("Uploading...")

	fileTransfer.Start <- true
	<-fileTransfer.Done

	if fileTransfer.Error != nil {
		panic(fileTransfer.Error)
	}
	fmt.Println("Done!")
}
