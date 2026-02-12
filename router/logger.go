package router

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dxloc/go-logger"
)

type HttpLogger struct {
	http.ResponseWriter
	status int
}

func (l *HttpLogger) Header() http.Header {
	return l.ResponseWriter.Header()
}

func (l *HttpLogger) Write(b []byte) (int, error) {
	return l.ResponseWriter.Write(b)
}

func (l *HttpLogger) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

func Logger(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := &HttpLogger{
				ResponseWriter: w,
				status:         http.StatusOK,
			}
			start := time.Now()
			next.ServeHTTP(l, r)
			duration := time.Since(start)
			b := strings.Builder{}
			if l.status >= 400 {
				b.WriteString(logger.ColorMagenta)
			} else if l.status >= 200 {
				b.WriteString(logger.ColorGreen)
			} else if l.status >= 100 {
				b.WriteString(logger.ColorBlue)
			}
			fmt.Fprintf(&b, "%d %s", l.status, http.StatusText(l.status))
			b.WriteString(logger.ColorReset)
			logger.Printf(logger.INFO,
				"%s %s %s %s %s %s\n",
				r.RemoteAddr,
				r.Method,
				r.RequestURI,
				name,
				duration,
				b.String(),
			)
		})
	}
}
