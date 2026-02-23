package router

import (
	"crypto/sha256"
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

const SumSize = sha256.Size

type UploadResp struct {
	Name string `json:"name"`
	Sum  string `json:"sum"`
}

func ioReadAndSum(r io.Reader, dataType int) ([]byte, string, error) {
	var data []byte
	var e error

	n := 512
	if dataType == TypeBase64 {
		n = base64.StdEncoding.EncodedLen(512)
	}
	b := make([]byte, 0, n)
	h := sha256.New()
	done := false

	// Read and sum chunks of data
	for {
		if done {
			break
		}
		if n, e = r.Read(b[:cap(b)]); e != nil {
			done = true
			if e == io.EOF {
				e = nil
			}
		}
		if e != nil || n == 0 {
			break
		}
		if dataType == TypeBase64 {
			if n, e = base64.StdEncoding.Decode(b, b[:n]); e != nil {
				done = true
				break
			}
		}
		data = append(data, b[:n]...)
		h.Write(b[:n])
	}
	if e != nil {
		return nil, "", e
	}
	if len(data) == 0 {
		return nil, "", nil
	}

	return data, base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func Upload(w http.ResponseWriter, r *http.Request) {
	var path strings.Builder
	var fname strings.Builder
	var data []byte
	var dType int
	var sum string
	var e error

	t := time.Now()
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/octet-stream":
		dType = TypeRaw
	case "application/base64":
		dType = TypeBase64
	default:
		api.ThrowMessage(http.StatusBadRequest, "invalid content type")
	}
	if data, sum, e = ioReadAndSum(r.Body, dType); e != nil {
		logger.Error(e, "read body ", contentType)
		api.ThrowError(http.StatusInternalServerError, e)
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
	writter.Response(http.StatusCreated, UploadResp{Name: fname.String(), Sum: sum})
}
