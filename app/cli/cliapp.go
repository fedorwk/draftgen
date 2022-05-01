package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/fedorwk/draftgen/generator"
	"github.com/fedorwk/draftgen/util"
	"github.com/fedorwk/go-util/data/delimiterdetector"
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
	templateFile, err := os.Open(TemplatePath)
	if err != nil {
		panic(err)
	}
	dataFile, err := os.Open(DataPath)
	if err != nil {
		panic(err)
	}

	if CSVDelim == "" {
		resolveDelimiter(dataFile)
	}

	items, err := util.ParseItems(dataFile, CSVDelim)
	if err != nil {
		panic(err)
	}

	generator := generator.DraftGenerator{
		Subject:          Subject,
		Items:            items,
		EmailPlaceholder: EmailPlaceholder,
		StartDelim:       StartDelim,
		EndDelim:         EndDelim,
	}
	err = generator.ParseTemplate(templateFile)
	if err != nil {
		panic(err)
	}
	// TODO:

	// generate dests files (io.Writer's)

	// execute generator to dests

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

func resolveDelimiter(src io.Reader) {
	var err error
	CSVDelim, err = delimiterdetector.Parse(src, 3)

	if err != nil {
		switch err {
		case delimiterdetector.ErrEmptySource:
			fmt.Println("Data source is empty\nExiting...")
			os.Exit(1)
		case delimiterdetector.ErrDefiningDelimiter:
			fmt.Println("Ragged CSV input\nExiting...")
			os.Exit(1)
		case delimiterdetector.ErrMultipleDelimiterOptions:
			fmt.Println("Can't detect the delimiter. please specify it")
			fmt.Scanln(&CSVDelim)
		}
	}
	if CSVDelim == "" {
		fmt.Println("Still unable to detect delimiter\nExiting...")
		os.Exit(1)
	}
}

// 	fmt.Println("Can't detect the delimiter. please specify it")
//	fmt.Scanln(&CSVDelim)
