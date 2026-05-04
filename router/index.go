package router

import (
	"cmp"
	"crypto/sha256"
	"encoding/base64"
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
	name  string
	date  string
	size  int
}

type FileEntryResp struct {
	Name    string `json:"name"`
	IsDir   bool   `json:"is_dir"`
	Date    string `json:"date"`
	Size    int    `json:"size"`
	Content string `json:"content,omitempty"`
	Sum     string `json:"sum,omitempty"`
}

type FolderEntryResp struct {
	Name    string          `json:"name"`
	IsDir   bool            `json:"is_dir"`
	Size    int             `json:"size"`
	Content []FileEntryResp `json:"content,omitempty"`
}

func commaFormat(i int) string {
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
				sizeInt := int(info.Size())
				size := commaFormat(sizeInt)
				if info.IsDir() {
					size = "-"
					sizeInt = 0
				}
				tz, _ := info.ModTime().Zone()
				date := fmt.Sprintf(
					"%04d-%02d-%02d %02d:%02d:%02d UTC%s",
					info.ModTime().Year(), info.ModTime().Month(), info.ModTime().Day(),
					info.ModTime().Hour(), info.ModTime().Minute(), info.ModTime().Second(),
					tz,
				)
				ie := IndexEntry{
					Name:  a.String(),
					IsDir: info.IsDir(),
					Size:  dom.NewDiv(size, dom.NewAttr("style", "text-align:right")),
					Date: dom.NewDiv(
						date,
						dom.NewAttr("style", "text-align:center"),
					),
					name: name.String(),
					size: sizeInt,
					date: date,
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
	f := r.URL.Query().Get("format")

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
			if f == "json" {
				b, e := os.ReadFile(path.String())
				if e != nil {
					logger.Error(e, "read file ", path)
					api.ThrowError(http.StatusInternalServerError, e)
				}
				ie := FileEntryResp{
					Name:    r.URL.Path,
					IsDir:   false,
					Size:    len(b),
					Content: base64.StdEncoding.EncodeToString(b),
					Sum:     fmt.Sprintf("%02x", sha256.Sum256(b)),
				}
				w := api.NewApiResponseWritter(w)
				w.Response(http.StatusOK, ie)
				return
			}
			http.ServeFile(w, r, path.String())
			return
		} else if r.URL.Path[len(r.URL.Path)-1] != '/' {
			api.ThrowMessage(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		}
	}

	ls := listFiles(r.URL.Path, path.String())
	if f == "json" {
		ie := FolderEntryResp{
			Name:  r.URL.Path,
			IsDir: true,
			Size:  len(ls),
		}
		for _, l := range ls {
			ie.Content = append(ie.Content, FileEntryResp{
				Name:  l.name,
				IsDir: l.IsDir,
				Size:  l.size,
				Date:  l.date,
			})
		}
		w := api.NewApiResponseWritter(w)
		w.Response(http.StatusOK, ie)
		return
	}

	doc := dom.NewDocument("Simple upload server", "", "")
	body := doc.Body()

	body.AddElement(dom.NewHeading(1, "Index of "+r.URL.Path))
	body.AddElement(dom.NewTable(ls, dom.NewAttr("style", "width:100%")))

	w.Write([]byte(doc.Serialize()))
}
