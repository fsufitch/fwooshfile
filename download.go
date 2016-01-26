package filebounce

import "net/http"
import "time"

func RegisterDownloadHandlers() {
	http.HandleFunc("/d/", handleDownloads)
}

func handleDownloads(w http.ResponseWriter, r *http.Request) {
	dlId := r.URL.Path[len("/d/"):]
	if len(dlId) > 0 {
		handleDownloaderPage(dlId, w, r)
	} else {
		http.NotFound(w, r)
	}
}

func handleDownloaderPage(dlId string, w http.ResponseWriter, r *http.Request) {
	writeStuff := func(dt DownloadTarget) {
		dt.StartFile("output.txt", "text/plain")
		dt.Stream <- []byte("hello\n")
		time.Sleep(1 * time.Second)
		dt.Stream <- []byte("world\n")
		time.Sleep(1 * time.Second)
		dt.Stream <- []byte("how are you today?\n")
		time.Sleep(1 * time.Second)
		dt.Stream <- []byte("this is download ID: " + dlId + "\n")
		close(dt.Stream)
	}

	dt := NewDownloadTarget(dlId, w)
	go dt.Download()
	go writeStuff(dt)

	_ = <-dt.Done
}

func tryFlushing(w http.ResponseWriter) {
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

type DownloadTarget struct {
	dlId string
	out http.ResponseWriter
	headersDone chan bool
	
	Stream chan []byte
	Done chan bool
}

func NewDownloadTarget(dlId string, w http.ResponseWriter) (dt DownloadTarget) {
	dt.dlId = dlId
	dt.out = w
	dt.headersDone = make(chan bool)
	dt.Stream = make(chan []byte)
	dt.Done = make(chan bool)
	return
}

func (dt DownloadTarget) Download() {
	_ = <-dt.headersDone
	for data := range dt.Stream  {
		dt.out.Write(data)
		tryFlushing(dt.out)
	}
	dt.Done <- true
}

func (dt DownloadTarget) StartFile(filename, mimetype string) {
	dt.out.Header().Set("Content-Type", mimetype)
	dt.out.Header().Set("Content-Disposition", "attachment; filename=" + filename)
	dt.out.WriteHeader(200)
	dt.headersDone <- true
}
