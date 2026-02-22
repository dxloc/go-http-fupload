package router

import (
	"cmp"
	"fmt"
	"go-http-fupload/api"
	"go-http-fupload/dom"
	"net/http"
	"os"
	"strings"

	"github.com/dxloc/go-logger"
	"github.com/dxloc/gosort"
)

type IndexEntry struct {
	Name  string `html:"Name"`
	IsDir bool
	Date  dom.HtmlElement `html:"Date"`
	Size  dom.HtmlElement `html:"Size (Bytes)"`
}

func commaFormat(i int64) string {
	s := fmt.Sprintf("%d", i)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

func (ie IndexEntry) Compare(a, b IndexEntry) int {
	if a.IsDir != b.IsDir {
		if a.IsDir {
			return -1
		} else {
			return 1
		}
	} else {
		return cmp.Compare(a.Name, b.Name)
	}
}

func listFiles(url, folderPath string) []IndexEntry {
	var ls []IndexEntry
	var name strings.Builder
	var fullPath strings.Builder

	if entries, e := os.ReadDir(folderPath); e != nil {
		logger.Error(e, "read dir ", folderPath)
		api.ThrowError(http.StatusInternalServerError, e)
	} else {
		for _, entry := range entries {
			if info, e := entry.Info(); e == nil {
				name.Reset()
				name.WriteString(entry.Name())
				fullPath.Reset()
				fullPath.WriteString(url)
				if info.IsDir() {
					name.WriteString("/")
				}
				fullPath.WriteString(name.String())
				a := dom.NewElement("a", name.String(), dom.NewHref(fullPath.String()))
				size := commaFormat(info.Size())
				if info.IsDir() {
					size = "-"
				}
				tz, _ := info.ModTime().Zone()
				ie := IndexEntry{
					Name:  a.String(),
					IsDir: info.IsDir(),
					Size:  dom.NewDiv(size, dom.NewAttr("style", "text-align:right")),
					Date: dom.NewDiv(
						fmt.Sprintf(
							"%04d-%02d-%02d %02d:%02d:%02d UTC%s",
							info.ModTime().Year(), info.ModTime().Month(), info.ModTime().Day(),
							info.ModTime().Hour(), info.ModTime().Minute(), info.ModTime().Second(),
							tz,
						),
						dom.NewAttr("style", "text-align:center"),
					),
				}
				ls = append(ls, ie)
			}
		}
	}
	s := gosort.NewSorter[IndexEntry](0)
	s.Sort(ls, 0)

	if url != "/" {
		name.Reset()
		fullPath.Reset()
		name.WriteString("../")
		ss := strings.Split(url, "/")
		for i := 0; i < len(ss)-2; i++ {
			if ss[i] != "" {
				fullPath.WriteString("/")
				fullPath.WriteString(ss[i])
			}
		}
		fullPath.WriteString("/")
		a := dom.NewElement("a", name.String(), dom.NewHref(fullPath.String()))
		ie := IndexEntry{
			Name: a.String(),
		}
		ls = append([]IndexEntry{ie}, ls...)
	}
	return ls
}

func Index(w http.ResponseWriter, r *http.Request) {
	var path strings.Builder

	uri := ProxyPass(r.URL.Path)
	path.WriteString(config.DownloadDir)

	for i := 1; i < len(uri); i++ {
		if uri[i-1] == '/' && uri[i] == '/' {
			continue
		}
		path.WriteByte(uri[i])
	}

	if info, e := os.Stat(path.String()); e != nil {
		logger.Error(e, "stat ", path)
		if os.IsNotExist(e) {
			api.ThrowError(http.StatusNotFound, e)
		}
	} else {
		if !info.IsDir() {
			http.ServeFile(w, r, path.String())
			return
		} else if r.URL.Path[len(r.URL.Path)-1] != '/' {
			api.ThrowMessage(http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
	}

	ls := listFiles(r.URL.Path, path.String())

	doc := dom.NewDocument("Simple upload server", "", "")
	body := doc.Body()

	body.AddElement(dom.NewHeading(1, "Index of "+r.URL.Path))
	body.AddElement(dom.NewTable(ls, dom.NewAttr("style", "width:100%")))

	w.Write([]byte(doc.Serialize()))
}
