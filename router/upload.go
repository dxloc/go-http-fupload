package router

import (
	"encoding/base64"
	"fmt"
	"go-http-fupload/api"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dxloc/go-logger"
)

const (
	TypeBase64 = iota + 1
	TypeRaw
)

func Upload(w http.ResponseWriter, r *http.Request) {
	var path strings.Builder
	var fname strings.Builder
	var data []byte
	var e error

	t := time.Now()
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/octet-stream":
		if data, e = io.ReadAll(r.Body); e != nil {
			logger.Error(e, "read body")
			api.ThrowError(http.StatusInternalServerError, e)
		}
	case "application/base64":
		if b64, e := io.ReadAll(r.Body); e != nil {
			logger.Error(e, "read body base64")
			api.ThrowError(http.StatusInternalServerError, e)
		} else {
			data = make([]byte, base64.URLEncoding.DecodedLen(len(b64)))
			if n, e := base64.URLEncoding.Decode(data, b64); e != nil {
				logger.Error(e, "decode base64")
				api.ThrowError(http.StatusBadRequest, e)
			} else if n == 0 {
				api.ThrowMessage(http.StatusBadRequest, "empty decoded data")
			}
		}
	default:
		api.ThrowMessage(http.StatusBadRequest, "invalid content type")
	}
	if len(data) == 0 {
		api.ThrowMessage(http.StatusBadRequest, "empty data")
	}

	fmt.Fprintf(&fname, "%04d-%02d-%02d_%02d-%02d-%02d-%06d_",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
		t.Nanosecond()/1000,
	)
	path.WriteString(config.UploadDir)
	path.WriteString("/")
	url := ProxyPass(r.URL.Path)
	ss := strings.Split(url, "/")
	for i, s := range ss {
		if s != "" {
			if i == len(ss)-1 {
				fname.WriteString(s)
			} else {
				path.WriteString(s)
				path.WriteString("/")
				if info, e := os.Stat(path.String()); e != nil {
					if os.IsNotExist(e) {
						if e = os.MkdirAll(path.String(), 0755); e != nil {
							logger.Error(e, "mkdir")
							api.ThrowError(http.StatusInternalServerError, e)
						}
					} else {
						logger.Error(e, "stat")
						api.ThrowError(http.StatusInternalServerError, e)
					}
				} else if !info.IsDir() {
					api.ThrowMessage(http.StatusBadRequest, "folder name already exists in the parent folder")
				}
			}
		} else if i == len(ss)-1 {
			api.ThrowMessage(http.StatusBadRequest, "file name must not end with '/'")
		}
	}

	path.WriteString(fname.String())
	pathStr := path.String()
	if _, e := os.Stat(pathStr); os.IsNotExist(e) {
		if e = os.WriteFile(pathStr, data, 0644); e != nil {
			logger.Error(e, "write file ", pathStr)
			api.ThrowError(http.StatusInternalServerError, e)
		}
	} else if e != nil {
		logger.Error(e, "stat ", pathStr)
		api.ThrowError(http.StatusInternalServerError, e)
	} else {
		api.ThrowMessage(http.StatusBadRequest, "file name already exists in the parent folder")
	}

	writter := api.NewApiResponseWritter(w)
	writter.Response(http.StatusCreated, nil, api.HttpHeader{
		Key:   "Location",
		Value: strings.Replace(pathStr, config.UploadDir, "", 1),
	})
}
