package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

// TransferClient encapsulates the client network connection
type TransferClient struct {
	BaseURL, Name, Key string
}

// FileTransfer represents one upload session
type FileTransfer struct {
	Path, Mimetype    string
	File              *os.File
	prepared, started bool
	Start, Done       chan bool
	Error             error
}

// CheckAlive verifies the status of the API
func (c TransferClient) CheckAlive() error {
	if c.BaseURL[len(c.BaseURL)-1] != "/"[0] {
		c.BaseURL += "/"
	}

	statusURL := c.BaseURL + "api/status"
	resp, err := http.Get(statusURL)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Non-OK status code: %d", resp.StatusCode)
	}

	statusBody := [2]byte{}

	_, _ = resp.Body.Read(statusBody[:])
	if !bytes.Equal(statusBody[:], []byte("OK")) {
		return fmt.Errorf("Non-OK query body: %s", string(statusBody[:]))
	}

	return nil
}

func (c TransferClient) registerNewUpload(filename, mimetype string, size int64) (uploadID string, err error) {
	url := c.BaseURL
	if url[len(url)-1] != "/"[0] {
		url += "/"
	}
	url += "api/new_upload/"

	req, err := http.NewRequest("POST", url, strings.NewReader(""))
	if err != nil {
		return
	}

	req.Header.Set("X-FileBounce-Filename", filename)
	req.Header.Set("X-FileBounce-Content-Type", mimetype)
	req.Header.Set("X-FileBounce-Content-Length", strconv.FormatInt(size, 10))

	req.Header.Write(os.Stderr)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		resp.Header.Write(os.Stderr)
		io.Copy(os.Stderr, resp.Body)
		fmt.Fprintln(os.Stderr, "")
		return "", fmt.Errorf("Non-200 status code: %d", resp.StatusCode)
	}

	bodyData, fileErr := ioutil.ReadAll(resp.Body)

	if fileErr != nil && fileErr != io.EOF {
		err = fileErr
		return
	}
	uploadID = string(bodyData)
	if uploadID == "" {
		return "", errors.New("No upload ID received!")
	}

	return
}

// Prepare prepares a file transfer and makes sure its state is ready for upload
// * check that File is non-nil, otherwise initialize it from Path
// * check that File can be opened and is not a directory
// * create the Done and Start channels
func (t *FileTransfer) Prepare() error {
	if t.prepared {
		return nil // already done
	}

	if len(t.Path) > 0 {
		if t.File != nil {
			return errors.New("Both path and separate file object specified")
		}
		file, err := os.Open(t.Path)
		if err != nil {
			return err
		}
		t.File = file
	} else {
		if t.File == nil {
			return errors.New("No path or explicit File specified")
		}
	}

	info, err := t.File.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		return errors.New("Given file must not be a directory")
	}

	if len(t.Mimetype) == 0 {
		ext := filepath.Ext(t.File.Name())
		t.Mimetype = mime.TypeByExtension(ext)
		if t.Mimetype == "" {
			t.Mimetype = "application/octet-stream"
		}
	}

	t.Done = make(chan bool)
	t.Start = make(chan bool)
	t.prepared = true
	return nil
}

// StartTransfer performs the upload using the given client
func (t *FileTransfer) StartTransfer(c TransferClient) (uploadID string, err error) {
	if err = c.CheckAlive(); err != nil {
		return
	}
	if err = t.Prepare(); err != nil {
		return
	}
	if t.started {
		return "", errors.New("This transfer already started")
	}

	info, _ := t.File.Stat()

	uploadID, err = c.registerNewUpload(
		filepath.Base(t.File.Name()),
		t.Mimetype,
		info.Size())
	if err != nil {
		return
	}

	t.started = true
	go t.doUploadWebSocket(uploadID, c)
	return
}

func (t *FileTransfer) doUploadWebSocket(uploadID string, c TransferClient) {
	wsURL := c.BaseURL + "api/upload_ws/" + uploadID
	if wsURL[:4] == "http" {
		wsURL = "ws" + wsURL[4:]
	} else {
		t.Error = errors.New("Unrecognized non-HTTP scheme: " + wsURL)
	}

	<-t.Start
	config, _ := websocket.NewConfig(wsURL, c.BaseURL)
	conn, err := websocket.DialConfig(config)
	//defer conn.Close()

	if err != nil {
		t.Error = err
		t.Done <- true
		return
	}

	buf := make([]byte, 50000) // 50 KB arbitrary
	reader := bufio.NewReader(t.File)
	for {
		_, fileErr := reader.Read(buf)
		if fileErr != nil {
			if fileErr == io.EOF {
				break
			}
			t.Error = fileErr
			t.Done <- true
			return
		}
		conn.Write(buf)
		conn.Read(buf)
		if !bytes.Equal(buf[:2], []byte("OK")) {
			t.Error = fmt.Errorf("Data transfer error! %s", string(buf))
			t.Done <- true
			return
		}
	}

	t.Done <- true
}
