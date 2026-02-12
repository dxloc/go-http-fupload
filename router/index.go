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

func listFiles(base, folderPath string) []IndexEntry {
	var ls []IndexEntry
	var name strings.Builder
	var fullPath strings.Builder

	if entries, e := os.ReadDir(base + folderPath); e != nil {
		logger.Error(e, "read dir ", base+folderPath)
		api.ThrowError(http.StatusInternalServerError, e)
	} else {
		for _, entry := range entries {
			if info, e := entry.Info(); e == nil {
				name.Reset()
				name.WriteString(entry.Name())
				fullPath.Reset()
				fullPath.WriteString(folderPath)
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

	if folderPath != "/" {
		name.Reset()
		fullPath.Reset()
		name.WriteString("../")
		ss := strings.Split(folderPath, "/")
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
	uri := r.URL.Path
	path := config.DownloadDir + uri

	if info, e := os.Stat(path); e != nil {
		logger.Error(e, "stat ", path)
		if os.IsNotExist(e) {
			api.ThrowError(http.StatusNotFound, e)
		}
	} else {
		if !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		} else if r.URL.Path[len(r.URL.Path)-1] != '/' {
			api.ThrowMessage(http.StatusBadRequest, "folder path must ends with '/'")
			return
		}
	}

	ls := listFiles(config.DownloadDir, uri)

	doc := dom.NewDocument("Simple upload server", "", "")
	body := doc.Body()

	body.AddElement(dom.NewHeading(1, "Index of "+uri))
	body.AddElement(dom.NewTable(ls, dom.NewAttr("style", "width:100%")))

	w.Write([]byte(doc.Serialize()))
}
