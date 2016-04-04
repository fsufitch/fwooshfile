package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fsufitch/fwooshfile/webui"
)

// Config of this bounce server
type Config struct {
	Host, StaticPath, TemplatePath string
	RelayAddrs                     []string
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

	fmt.Fprint(os.Stderr, "Setting up paths...\n")
	webui.InitializeWebUI(conf.StaticPath, conf.TemplatePath)

	fmt.Fprint(os.Stderr, "Setting up relays and round robin load balancer...\n")
	webui.Relays = []webui.FwooshRelay{}
	for _, relayAddr := range conf.RelayAddrs {
		webui.Relays = append(webui.Relays, webui.FwooshRelay{
			PrivateHost: relayAddr,
			PublicHost:  relayAddr,
		})
	}
	err = webui.RestartRoundRobinBalancer()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "Serving on %s.\n", conf.Host)
	err = http.ListenAndServe(conf.Host, nil)
	if err != nil {
		panic(err)
	}

}
