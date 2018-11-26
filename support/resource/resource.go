package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const maxUploadSize = 5 * 1024 // 2 mb
const uploadPath = "/Users/wjm/sourceTree/im_service/support/resource/res"

func main() {
	http.HandleFunc("/", handler())
	http.ListenAndServe(":8080", nil)
}

type UploadResp struct {
	Src string `json:"src"`
	Url string `json:"src_url"`
}

func handler() http.HandlerFunc {
	fileServer := http.FileServer(http.Dir(uploadPath))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "GET" {
			fileServer.ServeHTTP(w, r)
		} else if method == "POST" {

			fmt.Println(r)

			uri := r.RequestURI

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

			if uri == "/images" {
				fileType := r.Header.Get("File-Type")
				name := randToken(12) + "." + fileType

				newPath := filepath.Join(uploadPath, "images", name)

				// write file
				newFile, _ := os.Create(newPath)
				defer newFile.Close()

				if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
					renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				src := "/images/" + name

				bytes, _ := json.Marshal(UploadResp{Src: src, Url: "http://api.gobelieve.io" + src})
				w.Write(bytes)

				// 生成缩略图
				path, format, _ := isPictureFormat(newPath)
				compress(
					func() (io.Reader, error) {
						return os.Open(path)
					},
					func() (*os.File, error) {
						return os.Open(path)
					},
					path+"@256w_256h_0c", // 缩略图名称
					75,
					256, // 宽
					format)

			} else if uri == "/audios" {
				name := randToken(12)
				newPath := filepath.Join(uploadPath, "audios", name)
				// write file
				newFile, _ := os.Create(newPath)
				defer newFile.Close()
				if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
					renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				src := "/audios/" + name
				bytes, _ := json.Marshal(UploadResp{Src: src, Url: "http://api.gobelieve.io" + src})
				w.Write(bytes)
			}
		}
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

func compress(
	getReadSizeFile func() (io.Reader, error),
	getDecodeFile func() (*os.File, error),
	to string,
	quality,
	base int,
	format string) bool {

	/** 读取文件 */
	file_origin, err := getDecodeFile()
	defer file_origin.Close()
	if err != nil {
		log.Fatal(err)
		return false
	}
	var origin image.Image
	var config image.Config
	var temp io.Reader

	/** 读取尺寸 */
	temp, err = getReadSizeFile()
	if err != nil {
		log.Fatal(err)
		return false
	}
	var typeImage int64
	format = strings.ToLower(format)
	/** jpg 格式 */
	if format == "jpg" || format == "jpeg" {
		typeImage = 1
		origin, err = jpeg.Decode(file_origin)
		if err != nil {
			log.Fatal(err)
			return false
		}
		temp, err = getReadSizeFile()
		if err != nil {
			log.Fatal(err)
			return false
		}
		config, err = jpeg.DecodeConfig(temp)
		if err != nil {
			return false
		}
	} else if format == "png" {
		typeImage = 0
		origin, err = png.Decode(file_origin)
		if err != nil {
			log.Fatal(err)
			return false
		}
		temp, err = getReadSizeFile()
		if err != nil {
			log.Fatal(err)
			return false
		}
		config, err = png.DecodeConfig(temp);
		if err != nil {
			return false
		}
	}

	/** 做等比缩放 */
	width := uint(base) /** 基准 */
	height := uint(base * config.Height / config.Width)

	if strings.Index(to, "_h") > -1 {
		to = strings.Replace(to, "_h", fmt.Sprintf("_%dh", height), 1)
	}

	canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)
	out, err := os.Create(to)
	defer out.Close()
	if err != nil {
		log.Fatal(err)
		return false
	}
	if typeImage == 0 {
		err = png.Encode(out, canvas)
		if err != nil {
			fmt.Println("压缩图片失败")
			return false
		}
	} else {
		err = jpeg.Encode(out, canvas, &jpeg.Options{quality})
		if err != nil {
			fmt.Println("压缩图片失败")
			return false
		}
	}

	return true
}

/** 是否是图片 */
func isPictureFormat(path string) (string, string, string) {
	temp := strings.Split(path, ".")
	if len(temp) <= 1 {
		return "", "", ""
	}

	mapRule := make(map[string]int64)
	mapRule["jpg"] = 1
	mapRule["png"] = 1
	mapRule["jpeg"] = 1

	/** 添加其他格式 */
	if mapRule[temp[1]] == 1 {
		println(temp[1])
		return path, temp[1], temp[0]
	} else {
		return "", "", ""
	}
}

func Compress(path string, quality int, width int) {
	path, format, _ := isPictureFormat(path)
	//out := fmt.Sprintf("%s_@%dw_h.%s", parent, width, format)
	out := path + "@256w_256h_0c"
	if !compress(
		func() (io.Reader, error) {
			return os.Open(path)
		},
		func() (*os.File, error) {
			return os.Open(path)
		},
		out,
		quality,
		width,
		format) {
		fmt.Println("生成缩略图失败")
	} else {
		fmt.Println("生成缩略图成功")
	}
}
