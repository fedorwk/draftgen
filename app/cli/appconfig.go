package cli

import "time"

type config struct {
	LinesCountToAnalyzeCSV int

	DefauluStartDelim string
	DefaultEndDelim   string

	OutputDir        string
	OutputFileSuffix string
	OutputZipName    string
}

var appConfig = config{
	LinesCountToAnalyzeCSV: 3,

	DefauluStartDelim: "{",
	DefaultEndDelim:   "}",

	OutputDir:        "output_" + time.Now().Format(time.RFC822),
	OutputFileSuffix: ".eml",
	OutputZipName:    "output.zip",
}
