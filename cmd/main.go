package main

import (
	"flag"
	"go-http-fupload/router"
	"log"
	"net/http"

	"github.com/dxloc/go-logger"
	"gopkg.in/ini.v1"
)

// Default values
var configFile = ".env"
var baseUri = "/"
var port = "8080"
var uploadDir = "upload"
var downloadDir = "files"
var logLevel = "debug"

func init() {
	// Parse flags
	configFileParse := flag.String("config", configFile, "Config file")
	baseUriParse := flag.String("base-uri", baseUri, "Base URI")
	portParse := flag.String("port", port, "Port")
	uploadDirParse := flag.String("upload-dir", uploadDir, "Upload directory")
	downloadDirParse := flag.String("download-dir", downloadDir, "File directory")
	logLevelParse := flag.String("log-level", logLevel, "Log level")

	flag.Parse()

	// Set config file
	if *configFileParse != "" {
		configFile = *configFileParse
	}

	// Parse configFile file
	if env, err := ini.Load(configFile); err == nil {
		if v, e := env.Section("").GetKey("BASE_URI"); e == nil {
			*baseUriParse = v.String()
		}
		if v, e := env.Section("").GetKey("PORT"); e == nil {
			*portParse = v.String()
		}
		if v, e := env.Section("").GetKey("UPLOAD_DIR"); e == nil {
			*uploadDirParse = v.String()
		}
		if v, e := env.Section("").GetKey("FILE_DIR"); e == nil {
			*downloadDirParse = v.String()
		}
		if v, e := env.Section("").GetKey("LOG_LEVEL"); e == nil {
			*logLevelParse = v.String()
		}
	}

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
}

func main() {
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
