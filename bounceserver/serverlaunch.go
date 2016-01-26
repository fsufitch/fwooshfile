package main

import "net/http"

import "filebounce"

func main() {
	// TODO: dynamic config

	filebounce.RegisterDownloadHandlers()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
