package cli

import (
	"flag"
	"fmt"
	"os"
)

var (
	TemplatePath string
	DataPath     string

	EmailPlaceholder string
	Subject          string

	StartDelim string
	EndDelim   string

	CSVDelim string
)

func Run() {
	resolveFlags()
}

func resolveFlags() {
	TemplatePath = *flag.String("t", "", "template path")
	DataPath = *flag.String("d", "", "data path")

	EmailPlaceholder = *flag.String("e", "", "email placeholder")
	Subject = *flag.String("s", "", "email subject")

	StartDelim = *flag.String("a", "{", "placeholder start")
	EndDelim = *flag.String("b", "}", "placeholder end")

	CSVDelim = *flag.String("c", "", "csv delimiter")

	flag.Parse()

	if TemplatePath == "" {
		TemplatePath = flag.Arg(0)
		if TemplatePath == "" {
			fmt.Println(ErrNoTemplateArg)
			os.Exit(1)
		}
	}
	if DataPath == "" {
		DataPath = flag.Arg(0)
		if DataPath == "" {
			fmt.Println(ErrNoDataArg)
			os.Exit(1)
		}
	}
}
