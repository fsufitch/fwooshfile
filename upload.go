package filebounce

import "net/http"
import "strconv"

func RegisterUploadHandlers() {
  http.HandleFunc("/new_upload/", handleNewUpload)
  //http.HandleFunc("/upload/", handleUploadChunk)
}

func handleNewUpload(w http.ResponseWriter, r *http.Request) {
  filename := r.Header.Get("X-FileBounce-Filename")
  mimetype := r.Header.Get("X-FileBounce-Content-Type")
  size, size_err := strconv.Atoi(r.Header.Get("X-FileBounce-Content-Length"))
  token := r.Header.Get("X-FileBounce-Token")

  if len(filename) == 0 {
    http.Error(w, "No filename specified", 400)
    return
  }
  if len(mimetype) == 0 {
    http.Error(w, "No MIME type specified", 400)
    return
  }
  if size_err != nil {
    http.Error(w, "Invalid content length specified", 400)
    return
  }

  bf := NewBounceFile(filename, mimetype, token, size)

  http.SetCookie(w, &http.Cookie{
    Name: bf.CookieName,
    Value: bf.DlId,
  })
  w.Write([]byte(string(bf.DlId)))
}
