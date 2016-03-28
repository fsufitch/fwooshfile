package filebounce

import "net/http"

func RegisterStatusHandlers() {
	http.HandleFunc("/api/status", handleStatusPage)
}

func handleStatusPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}
