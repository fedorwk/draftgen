//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/fedorwk/draftgen/generator"
	"github.com/fedorwk/draftgen/util"
)

func main() {
	js.Global().Set("generate_go", js.FuncOf(Generate))
	js.Global().Set("retrieveResult_go", js.FuncOf(WriteBufferToJSUint8Array))

	<-make(chan bool)
}

var (
	OutputDataBuffer *bytes.Buffer
)

func Generate(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		fmt.Println(ErrRaggedInput)
		return -1
	}
	formData := args[0]
	generator, err := parseInput(formData)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	filenames := util.GenerateFilenames(generator.Items, config.FilenameGenFunc)
	OutputDataBuffer = &bytes.Buffer{}
	err = generator.Zip(OutputDataBuffer, filenames)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return OutputDataBuffer.Len()
}

func WriteBufferToJSUint8Array(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		fmt.Println("error: no dest passed")
		return -1
	}
	dst := args[0]
	nBytes := js.CopyBytesToJS(dst, OutputDataBuffer.Bytes())
	return nBytes
}

func parseInput(input js.Value) (*generator.DraftGenerator, error) {
	if !input.Truthy() {
		return nil, ErrNullInput
	}
	subject := input.Get("subject").String()
	template := input.Get("template").String()
	startDelim := input.Get("start_delim").String()
	endDelim := input.Get("end_delim").String()
	csvDelim := input.Get("csv_delim").String()

	// TODO: Can't get len of data
	jsInputData := input.Get("data")
	dataLen := jsInputData.Length()

	inputData := make([]byte, dataLen)
	js.CopyBytesToGo(inputData, input.Get("data"))
	inputDataReader := bytes.NewReader(inputData)
	items, _, err := util.ParseItems(inputDataReader, csvDelim)
	if err != nil {
		return nil, err
	}

	emailPlaceholder := util.DefineEmailPlaceholder(items)

	generator := &generator.DraftGenerator{
		Subject:          subject,
		Template:         template,
		Items:            items,
		EmailPlaceholder: emailPlaceholder,
		StartDelim:       startDelim,
		EndDelim:         endDelim,
	}
	return generator, nil
}

type Config struct {
	LinesCountToAnalyzeCSV int

	DefauluStartDelim string
	DefaultEndDelim   string

	OutputFileSuffix string
	OutputZipName    string

	FilenameGenFunc func(index int, item map[string]string) string
}

var config = Config{
	LinesCountToAnalyzeCSV: 3,

	DefauluStartDelim: "{",
	DefaultEndDelim:   "}",

	OutputFileSuffix: ".eml",
	OutputZipName:    "output.zip",

	FilenameGenFunc: func(index int, item map[string]string) string {
		return strconv.Itoa(index+1) + ".eml"
	},
}

var (
	ErrNullInput   = errors.New("err: null input passed")
	ErrRaggedInput = errors.New("err: ragged input, 2 passed values expected: fromData and destination")
)
