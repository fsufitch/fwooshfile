package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fsufitch/fwooshfile/relay"
)

// Config of this bounce server
type Config struct {
	Host string
}

func main() {
	conf := Config{":8080"}

	relay.RegisterDownloadHandlers()
	relay.RegisterUploadHandlers()
	relay.RegisterStatusHandlers()

	fmt.Fprintf(os.Stderr, "Serving on %s\n", conf.Host)
	err := http.ListenAndServe(conf.Host, nil)
	if err != nil {
		panic(err)
	}

}
