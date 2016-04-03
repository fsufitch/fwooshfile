package main

import (
	"encoding/json"
	"filebounce"
	"filebounce/webui"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Config of this bounce server
type Config struct {
	Host, StaticPath, TemplatePath string
}

func parseConfiguration(confpath string) (c Config, err error) {
	data, err := ioutil.ReadFile(confpath)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &c)
	return
}

func main() {
	confpath := flag.String("conf", "/etc/filebounce.json", "path to JSON config file")
	flag.Parse()

	conf, err := parseConfiguration(*confpath)
	if err != nil {
		panic(err)
	}

	filebounce.RegisterDownloadHandlers()
	filebounce.RegisterUploadHandlers()
	filebounce.RegisterStatusHandlers()

	fmt.Fprint(os.Stderr, "Setting up paths...\n")
	webui.InitializeWebUI(conf.StaticPath, conf.TemplatePath)

	fmt.Fprintf(os.Stderr, "Serving on %s\n", conf.Host)
	err = http.ListenAndServe(conf.Host, nil)
	if err != nil {
		panic(err)
	}

}
