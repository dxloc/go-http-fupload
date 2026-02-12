package api

import (
	"encoding/json"
	"net/http"
)

type HttpHeader struct {
	Key   string
	Value string
}

type ApiResponseWritter interface {
	Response(code int, body any, headers ...HttpHeader)
}

type ApiResponseWritterImpl struct {
	w http.ResponseWriter
}

func (i *ApiResponseWritterImpl) Response(code int, body any, headers ...HttpHeader) {
	i.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	for _, h := range headers {
		i.w.Header().Set(h.Key, h.Value)
	}
	i.w.WriteHeader(code)
	if body == nil {
		return
	}
	if b, e := json.Marshal(body); e == nil {
		i.w.Write(b)
	} else {
		e := ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: e.Error(),
		}
		b, _ := json.Marshal(e)
		i.w.WriteHeader(int(e.Status))
		i.w.Write(b)
	}
}

func NewApiResponseWritter(w http.ResponseWriter) ApiResponseWritter {
	return &ApiResponseWritterImpl{w}
}
