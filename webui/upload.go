package webui

import "net/http"

type uploadPageData struct {
	RelayURL string
}

func handleUploadPage(w http.ResponseWriter, r *http.Request) {
	select {
	case relayURL := <-RoundRobinBalancer:
		data := uploadPageData{relayURL.PublicHost}
		err := Templates.ExecuteTemplate(w, "upload", data)
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), 500)
		}
	default:
		http.Error(w, "No relay returned by load balancer!", 500)
	}
}
