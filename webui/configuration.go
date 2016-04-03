package webui

import "html/template"
import "net/http"
import "path"

// Templates is a html/template.Template encompassing available templates
var Templates *template.Template

// InitializeWebUI takes a static path and a template path, and does what it says on the can
func InitializeWebUI(staticpath, templatepath string) {
	_Templates, err := template.ParseGlob(path.Join(templatepath, "*"))
	Templates = _Templates
	if err != nil {
		panic(err)
	}
	println(staticpath)
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir(staticpath))))
	http.HandleFunc("/", handleUploadPage)

}
