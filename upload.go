package filebounce

import "fmt"
import "io/ioutil"
import "net/http"
import "strconv"

func RegisterUploadHandlers() {
  http.HandleFunc("/new_upload/", handleNewUpload)
  http.HandleFunc("/upload/", handleUploadChunk)
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

func handleUploadChunk(w http.ResponseWriter, r *http.Request) {
  dlId := r.URL.Path[len("/upload/"):]
  bf := GetBounceFile(dlId)
  if bf == nil {
    http.NotFound(w, r)
    return
  }

/*
  c, err := r.Cookie(bf.CookieName)
  if err != nil || (c != nil && c.Value != bf.DlId) {
    http.Error(w, fmt.Sprintf("Invalid cookie value: ", c.Value), 400)
    return
  }
  */

  data, err := ioutil.ReadAll(r.Body)
  if err != nil {
    http.Error(w, "Error reading request body: " + err.Error(), 500)
    return
  }

  err = bf.SendData(data)
  if err != nil {
    http.Error(w, "Error bouncing data: " + err.Error(), 500)
    return
  }
  fmt.Println("[upload] Sent data:", data)
  if bf.TransferFinished {
    w.Write([]byte("Done"))
  } else {
    w.Write([]byte("OK\n"))
  }
}
