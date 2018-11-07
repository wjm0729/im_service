package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const maxUploadSize = 5 * 1024 // 2 mb
const uploadPath = "/Users/wjm/sourceTree/im_service/support/resource/res"

// https://github.com/zupzup/golang-http-file-upload-download/blob/master/main.go
func main() {
	http.Handle("/", http.FileServer(http.Dir(uploadPath)))
	http.HandleFunc("/upload", uploadFileHandler())
	http.ListenAndServe(":8080", nil)
}

type UploadResp struct {
	Url string `json:"src_url"`
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("upload", r)

		http.MaxBytesReader(w, r.Body, maxUploadSize)

		contentLength := r.Header.Get("Content-Length")
		size, _ := strconv.ParseInt(contentLength, 10, len(contentLength))
		if size > maxUploadSize {
			fmt.Println("body no data")
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		}

		fileBytes, _ := ioutil.ReadAll(r.Body)
		if len(fileBytes) == 0 {
			fmt.Println("body no data")
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}

		fileType := r.Header.Get("File-Type")

		name := randToken(12) + "." + fileType

		newPath := filepath.Join(uploadPath+"/files", name)

		// write file
		newFile, _ := os.Create(newPath)
		defer newFile.Close()

		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		bytes, _ := json.Marshal(UploadResp{Url: "http://api.gobelieve.io/files/" + name})
		w.Write(bytes)
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
