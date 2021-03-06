package relay

import "net/http"
import "strconv"

func RegisterDownloadHandlers() {
	http.HandleFunc("/d/", handleDownloads)
}

func handleDownloads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dlId := r.URL.Path[len("/d/"):]
	if len(dlId) > 0 {
		handleActualDownload(dlId, w, r)
	} else {
		http.NotFound(w, r)
	}
}

func handleActualDownload(dlId string, w http.ResponseWriter, r *http.Request) {
	dt := NewDownloadTarget(dlId, "sample", w)
	err := RegisterDownloadTarget(dlId, dt)
	if err != nil {
		http.Error(w, "Error registering download: "+err.Error(), 500)
		return
	}
	go dt.Download()
	<-dt.Done
}

func tryFlushing(w http.ResponseWriter) {
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

type DownloadTarget struct {
	dlId, Name  string
	out         http.ResponseWriter
	headersDone chan bool

	Stream chan []byte
	Done   chan bool
}

func NewDownloadTarget(dlId, name string, w http.ResponseWriter) (dt DownloadTarget) {
	dt.dlId = dlId
	dt.Name = name
	dt.out = w
	dt.headersDone = make(chan bool)
	dt.Stream = make(chan []byte)
	dt.Done = make(chan bool)
	return
}

func (dt DownloadTarget) Download() {
	//<-dt.headersDone
	for data := range dt.Stream {
		//fmt.Printf("[download] Received data (%d)\n", len(data))
		dt.out.Write(data)
		tryFlushing(dt.out)
		//fmt.Println("[download] Processed.")
	}
	dt.Done <- true
}

func (dt DownloadTarget) StartFile(bf *BounceFile) {
	dt.out.Header().Set("Content-Type", bf.Mimetype)
	dt.out.Header().Set("Content-Disposition", "attachment; filename="+bf.Filename)
	dt.out.Header().Set("Content-Length", strconv.Itoa(bf.Size))
	dt.out.WriteHeader(200)
	tryFlushing(dt.out)
	//dt.headersDone <- true
}
