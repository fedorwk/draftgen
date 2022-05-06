package server

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fedorwk/draftgen/util"
	"github.com/pkg/errors"
)

type Config struct {
	Server    *ServerConfig
	App       *AppConfig
	HTTPForm  *HTTPFormConfig
	Generator *GeneratorConfig
}

// Defaults
var config = Config{
	Server: &ServerConfig{
		Port: "8080",
	},
	App: &AppConfig{
		EmbeddedHTMLURL:   "https://github.com/fedorwk/homepage/raw/main/services/draftgen/embedded.html",
		EmbeddedHTMLLocal: "",
		FilenameGenFunc: func(index int, item map[string]string) string {
			return strconv.Itoa(index+1) + ".eml"
		},
		OutputFilename: "drafts.zip",
	},
	HTTPForm: &HTTPFormConfig{
		Template: "template",
		Data:     "data_file",

		EmailPlaceholder: "email_ph",
		Subject:          "subject",

		StartDelim: "start_delim",
		EndDelim:   "end_delim",

		CSVDelim: "data_delim",
	},
	Generator: &GeneratorConfig{
		LinesCountToAnalyzeCSV: 3,
	},
}

// ServerConfig contains values for http.Server and Endpoints
type ServerConfig struct {
	Port string
}
type AppConfig struct {
	// Load service front-end from remote source.
	// Remote source has priority over local source
	EmbeddedHTMLURL string
	// Load service front-end from local file
	EmbeddedHTMLLocal string

	// Specifies how file names will be generated in the zip archive
	// Look util/util_test.go ExampleGenerateFilenames for more info
	FilenameGenFunc util.NameGenFn
	// the name of the archive that will be sent to the user
	OutputFilename string
}

// HTTPFormConfig specifies value names in HTTP POST request
type HTTPFormConfig struct {
	Template string
	Data     string

	EmailPlaceholder string
	Subject          string

	StartDelim string
	EndDelim   string

	CSVDelim string
}

// GeneratorConfig contains values related to draftgen.Generator
type GeneratorConfig struct {
	LinesCountToAnalyzeCSV int
}

func init() {
	// TODO: READ CONFIG or ENV

	// Load Inner HTML of service into memory
	if config.App.EmbeddedHTMLURL != "" {
		EmbeddedHTML = bytesFromRemote(config.App.EmbeddedHTMLURL)
	} else if config.App.EmbeddedHTMLLocal != "" {
		EmbeddedHTML = bytesFromLocal(config.App.EmbeddedHTMLLocal)
	} else {
		log.Println("no source for service HTML given")
	}
}

func bytesFromRemote(url string) []byte {
	var res bytes.Buffer

	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(errors.Wrapf(err,
			"initializing server: can't get service HTML from remote source: %s",
			config.App.EmbeddedHTMLURL))
	}
	defer response.Body.Close()
	_, err = res.ReadFrom(response.Body)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "initializing server: err reading service HTML from response"))
	}
	return res.Bytes()
}

func bytesFromLocal(path string) []byte {
	var res bytes.Buffer

	file, err := os.Open(config.App.EmbeddedHTMLLocal)
	if err != nil {
		log.Fatalln(errors.Wrapf(err,
			"initializing server: can't get service HTML from local source: %s",
			config.App.EmbeddedHTMLLocal))
	}
	defer file.Close()
	_, err = res.ReadFrom(file)
	if err != nil {
		log.Fatalln(errors.Wrap(err,
			"initializing server: err reading service HTML form local file"))
	}
	return res.Bytes()
}
