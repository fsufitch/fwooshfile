package webui

import "html/template"
import "net/http"
import "path"

var Templates *template.Template

func InitializeWebUI(staticpath, templatepath string) {
  _Templates, err := template.ParseGlob(path.Join(templatepath, "*"))
  Templates = _Templates
  if err != nil { panic(err) }
  println(staticpath)
  http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir(staticpath))))
  http.HandleFunc("/", handleUploadPage)


}
