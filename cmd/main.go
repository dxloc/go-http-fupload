package cmd

import (
	"fmt"
	"go-http-fupload/getlongopts"
	"go-http-fupload/router"
	"log"
	"net/http"
	"os"

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
var Version = ""

func Init() {
	// Parse flags
	p := getlongopts.NewParser([]getlongopts.LongOption{
		{
			Long:    "config",
			Short:   "c",
			ArgType: getlongopts.ArgTypeFile, SetFn: func(value string) { configFile = value },
			Description: "Set configuration file",
		},
		{
			Long:    "base-uri",
			Short:   "b",
			ArgType: getlongopts.ArgTypeOther, SetFn: func(value string) { baseUri = value },
			Description: "Set base URI",
		},
		{
			Long:    "port",
			Short:   "p",
			ArgType: getlongopts.ArgTypeOther, SetFn: func(value string) { port = value },
			Description: "Set port",
		},
		{
			Long:    "upload-dir",
			Short:   "u",
			ArgType: getlongopts.ArgTypeDir, SetFn: func(value string) { uploadDir = value },
			Description: "Set upload directory",
		},
		{
			Long:    "download-dir",
			Short:   "d",
			ArgType: getlongopts.ArgTypeDir, SetFn: func(value string) { downloadDir = value },
			Description: "Set download directory",
		},
		{
			Long:    "log-level",
			Short:   "l",
			ArgType: getlongopts.ArgTypeOther, SetFn: func(value string) { logLevel = value },
			Description: "Set log level",
		},
		{
			Long:    "version",
			Short:   "v",
			ArgType: getlongopts.ArgTypeNone, SetFn: func(value string) { fmt.Println(Version); os.Exit(0) },
			Description: "Print version and exit",
		},
	})
	if _, e := p.Parse(os.Args); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}

	// Parse configFile file
	if env, err := ini.Load(configFile); err == nil {
		if v, e := env.Section("").GetKey("BASE_URI"); e == nil && v.String() != "" {
			baseUri = v.String()
		}
		if v, e := env.Section("").GetKey("PORT"); e == nil && v.String() != "" {
			port = v.String()
		}
		if v, e := env.Section("").GetKey("UPLOAD_DIR"); e == nil && v.String() != "" {
			uploadDir = v.String()
		}
		if v, e := env.Section("").GetKey("FILE_DIR"); e == nil && v.String() != "" {
			downloadDir = v.String()
		}
		if v, e := env.Section("").GetKey("LOG_LEVEL"); e == nil && v.String() != "" {
			logLevel = v.String()
		}
	}

	logger.InitStr(logLevel)
}

func Main() {

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
