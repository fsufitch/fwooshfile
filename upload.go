package filebounce

import "encoding/base64"
import "fmt"
import "io/ioutil"
import "net/http"
import "strconv"

func RegisterUploadHandlers() {
	http.HandleFunc("/api/new_upload/", handleNewUpload)
	http.HandleFunc("/api/upload/", handleUploadChunk)
	RegisterWSUploadHandlers()
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
  dlId := r.URL.Path[len("/api/upload/"):]
  bf := GetBounceFile(dlId)
  if bf == nil {
    http.NotFound(w, r)
    return
  }

  encodedData, err := ioutil.ReadAll(r.Body)
  if err != nil {
    http.Error(w, "Error reading request body: " + err.Error(), 500)
    return
  }

  data, err := base64.StdEncoding.DecodeString(string(encodedData))
  if err != nil {
    http.Error(w, "Error decoding base64 data: " + err.Error(), 500)
    return
  }

  err = bf.SendData(data)
  if err != nil {
    http.Error(w, "Error bouncing data: " + err.Error(), 500)
    return
  }
  if bf.TransferFinished {
    w.Write([]byte("Done\n"))
  } else {
    w.Write([]byte("OK\n"))
  }
  fmt.Print("");
}
