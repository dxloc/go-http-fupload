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

func UrlValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, c := range r.URL.Path {
			if c >= 'a' && c <= 'z' {
				continue
			}
			if c >= 'A' && c <= 'Z' {
				continue
			}
			if c >= '0' && c <= '9' {
				continue
			}
			if c == '/' || c == '-' || c == '_' || c == '.' {
				continue
			}
			api.ThrowMessage(http.StatusBadRequest, "invalid URL")
		}
		next.ServeHTTP(w, r)
	})

}

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
		handler = UrlValidator(handler)
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
