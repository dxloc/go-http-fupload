package dom

import (
	"os"

	"github.com/dxloc/go-logger"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

type MinifyType string

func (t MinifyType) String() string {
	return string(t)
}

const (
	TypeTextCss        MinifyType = "text/css"
	TypeTextJavascript MinifyType = "text/javascript"
	TypeTextHtml       MinifyType = "text/html"
)

const (
	cssCachePrefix  string = "css:"
	jsCachePrefix   string = "js:"
	htmlCachePrefix string = "html:"
)

var minifiedFiles map[string]string
var m *minify.M

func CanMinify() bool {
	return m != nil
}

func Minify(typ MinifyType, text string) (string, error) {
	return m.String(typ.String(), text)
}

func LoadCssRaw(f, s string) {
	if s, e := Minify(TypeTextCss, s); e != nil {
		logger.Error(e, "minify css ", f)
	} else {
		minifiedFiles[jsCachePrefix+f] = s
	}
}

func LoadCss(files ...string) {
	for _, f := range files {
		if b, e := os.ReadFile(f); e != nil {
			logger.Error(e, "read file ", f)
		} else {
			LoadCssRaw(f, string(b))
		}
	}
}

func GetCss(file string) (string, bool) {
	s, f := minifiedFiles[cssCachePrefix+file]
	return s, f
}

func LoadJsRaw(f, s string) {
	if s, e := Minify(TypeTextJavascript, s); e != nil {
		logger.Error(e, "minify javascript ", f)
	} else {
		minifiedFiles[jsCachePrefix+f] = s
	}
}

func LoadJs(files ...string) {
	for _, f := range files {
		if b, e := os.ReadFile(f); e != nil {
			logger.Error(e, "read file ", f)
		} else {
			LoadJsRaw(f, string(b))
		}
	}
}

func GetJs(file string) (string, bool) {
	s, f := minifiedFiles[jsCachePrefix+file]
	return s, f
}

func LoadHtmlRaw(f, s string) {
	if s, e := Minify(TypeTextHtml, s); e != nil {
		logger.Error(e, "minify html ", f)
	} else {
		minifiedFiles[htmlCachePrefix+f] = s
	}
}

func LoadHtml(files ...string) {
	for _, f := range files {
		if b, e := os.ReadFile(f); e != nil {
			logger.Error(e, "read file ", f)
		} else {
			LoadHtmlRaw(f, string(b))
		}
	}
}

func GetHtml(file string) (string, bool) {
	s, f := minifiedFiles[htmlCachePrefix+file]
	return s, f
}

func MinifiedCacheInit(addr string) {
	minifiedFiles = make(map[string]string)
	m = minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("text/html", html.Minify)
}
