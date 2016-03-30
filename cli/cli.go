package main

import (
	"filebounce/client"
	"flag"
	"fmt"
	"os"
)

func main() {
	apiurlPtr := flag.String("apiurl", "http://localhost:8080",
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

	_, err := fileTransfer.StartTransfer(transferClient)

	if err != nil {
		panic(err)
	}
	<-fileTransfer.Done
	fmt.Println("Done!")
}
