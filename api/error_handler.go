package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dxloc/go-logger"
)

type ErrorResponseWriterImpl struct {
	w http.ResponseWriter
}

func (i *ErrorResponseWriterImpl) Response(code int, body any, headers ...HttpHeader) {
	i.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	for _, h := range headers {
		i.w.Header().Set(h.Key, h.Value)
	}
	b, _ := json.Marshal(body)
	i.w.WriteHeader(code)
	i.w.Write(b)
}

func NewErrorResponseWriter(w http.ResponseWriter) ApiResponseWritter {
	return &ErrorResponseWriterImpl{w}
}

func ThrowMessage(code int, message string) {
	e := fmt.Errorf("%03d:%s", code, message)
	panic(e)
}

func ThrowError(code int, e error) {
	ThrowMessage(code, e.Error())
}

func PanicHandler(w http.ResponseWriter, r *http.Request) {
	e := recover()
	if e == nil {
		return
	}

	var msg strings.Builder
	var code int

	writer := NewErrorResponseWriter(w)
	str := fmt.Sprint(e)
	logger.Error(errors.New(str), "api panic recovered")
	for i, s := range str {
		if s == ':' {
			if code, e = strconv.Atoi(str[:i]); e != nil {
				code = http.StatusInternalServerError
				msg.WriteString(str)
			} else {
				msg.WriteString(str[i+1:])
			}
			break
		}
	}
	if code == 0 {
		code = http.StatusInternalServerError
		msg.WriteString(str)
	}

	res := ApiResponse{Status: code, Message: msg.String()}
	writer.Response(code, res)
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer PanicHandler(w, r)
		next.ServeHTTP(w, r)
	})
}
