package main

import (
	"flag"
	"go-http-fupload/router"
	"log"
	"net/http"

	"github.com/dxloc/go-logger"
	"gopkg.in/ini.v1"
)

func main() {
	// Default values
	baseUri := "/"
	port := "8080"
	uploadDir := "upload"
	downloadDir := "files"
	logLevel := "debug"

	// Parse .env file
	if env, err := ini.Load(".env"); err == nil {
		if v, e := env.Section("").GetKey("BASE_URI"); e == nil {
			baseUri = v.String()
		}
		if v, e := env.Section("").GetKey("PORT"); e == nil {
			port = v.String()
		}
		if v, e := env.Section("").GetKey("UPLOAD_DIR"); e == nil {
			uploadDir = v.String()
		}
		if v, e := env.Section("").GetKey("FILE_DIR"); e == nil {
			downloadDir = v.String()
		}
		if v, e := env.Section("").GetKey("LOG_LEVEL"); e == nil {
			logLevel = v.String()
		}
	}

	// Parse flags
	baseUriParse := flag.String("base-uri", "", "Base URI")
	portParse := flag.String("port", "", "Port")
	uploadDirParse := flag.String("upload-dir", "", "Upload directory")
	downloadDirParse := flag.String("download-dir", "", "File directory")
	logLevelParse := flag.String("log-level", "", "Log level")

	flag.Parse()

	if *baseUriParse != "" {
		baseUri = *baseUriParse
	}
	if *portParse != "" {
		port = *portParse
	}
	if *uploadDirParse != "" {
		uploadDir = *uploadDirParse
	}
	if *downloadDirParse != "" {
		downloadDir = *downloadDirParse
	}
	if *logLevelParse != "" {
		logLevel = *logLevelParse
	}

	logger.InitStr(logLevel)

	log.Printf("Server started: %s:%s%s", "http://localhost", port, baseUri)
	log.Printf("Upload directory: %s", uploadDir)
	log.Printf("Download directory: %s", downloadDir)
	log.Printf("Log level: %s", logLevel)

	r := router.NewRouter(router.Config{
		BaseUri:     baseUri,
		Title:       "Simple upload server",
		UploadDir:   uploadDir,
		DownloadDir: downloadDir,
	})
	log.Fatal(http.ListenAndServe(":"+port, r))
}
