package main

import "net/http"
import "os"

import "filebounce"
import "filebounce/webui"

// Config of this bounce server
type Config struct {
	host, staticpath, templatepath string
}

func main() {
	// TODO: dynamic config
	staticpath := os.Args[1]
	templatepath := os.Args[2]

	filebounce.RegisterDownloadHandlers()
	filebounce.RegisterUploadHandlers()
	filebounce.RegisterStatusHandlers()
	webui.InitializeWebUI(staticpath, templatepath)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
