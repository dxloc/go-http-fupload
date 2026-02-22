package router

import (
	"go-http-fupload/api"
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	HandlerFunc http.HandlerFunc
}

type Config struct {
	BaseUri     string
	Title       string
	UploadDir   string
	DownloadDir string
}

var config Config

func ProxyPass(url string) string {
	if len(url) < len(config.BaseUri) {
		api.ThrowMessage(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}
	if url[:len(config.BaseUri)] == config.BaseUri {
		if len(url) == len(config.BaseUri) {
			url = "/"
		} else {
			url = "/" + url[len(config.BaseUri):]
		}
	} else {
		api.ThrowMessage(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	return url
}

func NewRouter(cfg Config) *mux.Router {
	var handler http.Handler

	config = cfg
	router := mux.NewRouter()
	routes := buildRoutes()

	for _, route := range routes {
		handler = route.HandlerFunc
		handler = api.PanicMiddleware(handler)
		handler = Logger(route.Name)(handler)
		router.
			Methods(route.Method).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func buildRoutes() []route {
	return []route{
		{"Index", "GET", Index},
		{"Upload", "POST", Upload},
	}
}
