package server

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fedorwk/draftgen/util"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	App      AppConfig      `yaml:"app"`
	HTTPForm HTTPFormConfig `yaml:"http_form"`
}

var config Config

// ServerConfig contains values for http.Server and Endpoints
type ServerConfig struct {
	Port string `yaml:"port" env:"PORT" env-default:"8080"`
}
type AppConfig struct {
	// Load service front-end from remote source.
	// Remote source has priority over local source
	InnerHTMLURL string `yaml:"inner_html_url" env:"INNER_HTML_URL"`
	// Load service front-end from local file
	InnerHTMLLocal string `yaml:"inner_html_local" env:"INNER_HTML_LOCAL"`

	// Specifies how file names will be generated in the zip archive
	// Look util/util_test.go ExampleGenerateFilenames for more info
	FilenameGenFunc util.NameGenFn
	// the name of the archive that will be sent to the user
	OutputFilename string `yaml:"output_filename" env:"OUTPUT_FNAME" env-default:"drafts.zip"`

	// Lines count will be parsed to detect csv delimiter, -1 for parse until EOF
	LinesCountToAnalyzeCSV int `yaml:"lines_to_analyze" env:"CSVNLINES" env-default:"3"`
}

// HTTPFormConfig specifies value names in HTTP POST request
type HTTPFormConfig struct {
	Template string `yaml:"template_http" env:"HTTPTEMPLATE" env-default:"template"`
	Data     string `yaml:"data_file_http" env:"HTTPDATAFILE" env-default:"data_file"`

	EmailPlaceholder string `yaml:"email_ph_http" env:"HTTPEMAILPH" env-default:"email_ph"`
	Subject          string `yaml:"subject_http" env:"HTTPSUBJECT" env-default:"subject"`

	StartDelim string `yaml:"start_delim_http" env:"HTTPSTARTDELIM" env-default:"start_delim"`
	EndDelim   string `yaml:"end_delim_http" env:"HTTPENDDELIM" env-default:"end_delim"`

	CSVDelim string `yaml:"data_delim_http" env:"HTTPCSVDELIM" env-default:"data_delim"`
}

func (c *Config) Init() error {
	var cfgpath string
	flagSet := flag.NewFlagSet("ServerMain", flag.ContinueOnError)
	flagSet.StringVar(&cfgpath, "cfg", "config.yml", "path to config file")
	flagSet.Usage = cleanenv.FUsage(flagSet.Output(), c, nil, flagSet.Usage)

	flagSet.Parse(os.Args[1:])
	err := cleanenv.ReadConfig(cfgpath, &config)
	if err != nil && err != io.EOF {
		return err
	}

	if config.App.InnerHTMLURL != "" {
		InnerHTML = bytesFromRemote(config.App.InnerHTMLURL)
	} else if config.App.InnerHTMLLocal != "" {
		InnerHTML = bytesFromLocal(config.App.InnerHTMLLocal)
	} else {
		log.Println("no source for service HTML given")
	}
	return nil
}

func bytesFromRemote(url string) []byte {
	var res bytes.Buffer

	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(errors.Wrapf(err,
			"initializing server: can't get service HTML from remote source: %s",
			config.App.InnerHTMLURL))
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

	file, err := os.Open(config.App.InnerHTMLLocal)
	if err != nil {
		log.Fatalln(errors.Wrapf(err,
			"initializing server: can't get service HTML from local source: %s",
			config.App.InnerHTMLLocal))
	}
	defer file.Close()
	_, err = res.ReadFrom(file)
	if err != nil {
		log.Fatalln(errors.Wrap(err,
			"initializing server: err reading service HTML form local file"))
	}
	return res.Bytes()
}
