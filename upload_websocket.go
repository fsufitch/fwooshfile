package filebounce

import "io"
import "net/http"
import "golang.org/x/net/websocket"

func RegisterWSUploadHandlers() {
	http.Handle("/api/upload_ws/", websocket.Handler(websocketUploadServer))
}

func websocketUploadServer(ws *websocket.Conn) {
	dlId := ws.Request().URL.Path[len("/api/upload_ws/"):]
	bf := GetBounceFile(dlId)
	if bf == nil {
		ws.Write([]byte("ERROR! Download with this ID does not exist."))
		return
	}

	var dataBuffer [50000]byte; // 50 KB arbitrary
	total := 0
	for {
		numBytes, err := ws.Read(dataBuffer[:])
		total += numBytes
		dataCopy := append([]byte{}, dataBuffer[:numBytes]...)
		
		bf.SendData(dataCopy)
		if err == io.EOF { break }
		
		if err != nil { // Oops!
			ws.Write([]byte("ERROR! " + err.Error()))
			break
		}
		ws.Write([]byte("OK"))
	}
}
