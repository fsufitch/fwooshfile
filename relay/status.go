package relay

import "net/http"

func RegisterStatusHandlers() {
	http.HandleFunc("/api/status", handleStatusPage)
}

func handleStatusPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}
