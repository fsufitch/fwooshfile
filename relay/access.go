package relay

import "net/http"

func ApplyCORS(w http.ResponseWriter, r *http.Request) (preflight, allow bool, err error) {
	allow = true // TODO: add actual CORS access control; this is a stub
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

	if r.Method == "OPTIONS" {
		preflight = true
	}

	return
}
