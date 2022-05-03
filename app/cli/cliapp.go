package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"github.com/fedorwk/draftgen/generator"
	"github.com/fedorwk/draftgen/util"
	"github.com/fedorwk/go-util/data/delimiterdetector"
)

var (
	TemplatePath  string
	DataPath      string
	OutputDirPath string
	ZipOutput     bool

	EmailPlaceholder string
	Subject          string

	StartDelim string
	EndDelim   string

	CSVDelim string
)

func Run(args []string) error {
	// DATA READING AND VALIDATION
	parseCliArgs(args)
	templateFile, err := os.Open(TemplatePath)
	if err != nil {
		return err
	}
	defer templateFile.Close()
	dataFile, err := os.Open(DataPath)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	// reason: file will be read twice if user did not specify delimiter
	var dataFileReader io.ReadWriter

	if CSVDelim == "" { // if delimiter not specified by user
		var streamCopy bytes.Buffer
		reader := io.TeeReader(dataFile, &streamCopy) // "clone" stream
		resolveDelimiter(reader)                      // read from cloned stream
		dataFileReader = &streamCopy
	} else {
		dataFileReader = dataFile
	}

	items, _, err := util.ParseItems(dataFileReader, CSVDelim)
	if err != nil {
		return err
	}

	// SETTING UP GENERATOR
	generator := &generator.DraftGenerator{
		Subject:          Subject,
		Items:            items,
		EmailPlaceholder: EmailPlaceholder,
		StartDelim:       StartDelim,
		EndDelim:         EndDelim,
	}
	err = generator.ParseTemplate(templateFile)
	if err != nil {
		return err
	}

	// WRITING OUTPUT
	err = makeOutputDir()
	if err != nil {
		return err
	}
	if ZipOutput {
		err = generateToZip(generator)
	} else {
		err = generateToFiles(generator)
	}
	if err != nil {
		return err
	}
	return nil
}

func parseCliArgs(args []string) {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&TemplatePath, "template", "", "template path")
	flags.StringVar(&TemplatePath, "t", "", "template path")
	flags.StringVar(&DataPath, "d", "", "data path")
	flags.StringVar(&DataPath, "data", "", "data path")
	flags.StringVar(&OutputDirPath, "output", appConfig.OutputDir, "generated drafts output dir path")
	flags.StringVar(&OutputDirPath, "o", appConfig.OutputDir, "generated drafts output dir path")
	flags.BoolVar(&ZipOutput, "zip", false, "zip output files")
	flags.BoolVar(&ZipOutput, "z", false, "zip output files")

	flags.StringVar(&EmailPlaceholder, "email", "", "email placeholder")
	flags.StringVar(&EmailPlaceholder, "e", "", "email placeholder")
	flags.StringVar(&Subject, "subject", "", "email subject")
	flags.StringVar(&Subject, "s", "", "email subject")

	flags.StringVar(&StartDelim, "sdelim", appConfig.DefauluStartDelim, "placeholder start")
	flags.StringVar(&StartDelim, "a", appConfig.DefauluStartDelim, "placeholder start")
	flags.StringVar(&EndDelim, "edelim", appConfig.DefaultEndDelim, "placeholder end")
	flags.StringVar(&EndDelim, "b", appConfig.DefaultEndDelim, "placeholder end")

	flags.StringVar(&CSVDelim, "delimiter", "", "csv delimiter")
	flags.StringVar(&CSVDelim, "c", "", "csv delimiter")

	flags.Parse(args)

	if TemplatePath == "" {
		TemplatePath = flags.Arg(0)
		if TemplatePath == "" {
			fmt.Println(ErrNoTemplateArg)
			os.Exit(1)
		}
	}
	if DataPath == "" {
		DataPath = flags.Arg(1)
		if DataPath == "" {
			fmt.Println(ErrNoDataArg)
			os.Exit(1)
		}
	}
}

// Wrap over delimiterdetector which returns user errors
func resolveDelimiter(src io.Reader) {
	var err error
	CSVDelim, err = delimiterdetector.Parse(src, appConfig.LinesCountToAnalyzeCSV)

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

func generateToFiles(gen *generator.DraftGenerator) error {
	dests, err := generateDestFiles(len(gen.Items))
	if err != nil {
		return err
	}
	for i, dst := range dests {
		err := gen.Execute(i, dst)
		if err != nil {
			return err
		}
		err = dst.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func generateDestFiles(n int) ([]*os.File, error) {
	err := os.Mkdir(OutputDirPath, 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return nil, err
	}
	dsts := make([]*os.File, 0, n)
	for i := 0; i < n; i++ {
		var dstfilePath strings.Builder
		dstfilePath.WriteString(OutputDirPath)
		dstfilePath.WriteString("/")
		dstfilePath.WriteString(strconv.Itoa(i + 1))
		dstfilePath.WriteString(appConfig.OutputFileSuffix)
		dstfile, err := os.Create(dstfilePath.String())
		if err != nil {
			return nil, err
		}
		dsts = append(dsts, dstfile)
	}

	return dsts, nil
}

func generateToZip(gen *generator.DraftGenerator) error {
	zipFile, err := os.Create(OutputDirPath + "/" + appConfig.OutputZipName)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	filenames := util.GenerateFilenames(gen.Items, func(index int, item map[string]string) string {
		return strconv.Itoa(index+1) + appConfig.OutputFileSuffix
	})
	err = gen.Zip(zipFile, filenames)
	if err != nil {
		return err
	}
	return nil
}

func makeOutputDir() error {
	err := os.Mkdir(OutputDirPath, 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	return nil
}
