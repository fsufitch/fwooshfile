package webui

import "net/http"

func handleUploadPage(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "upload", nil)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), 500)
	}
}
