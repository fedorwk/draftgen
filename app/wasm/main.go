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

func Generate(this js.Value, args []js.Value) interface{} {
	fmt.Println("DEBUG: generate_go() passed arguments", args)
	if len(args) != 2 {
		panic(ErrRaggedInput)
	}
	formData, dstJsObject := args[0], args[1]
	generator, err := parseInput(formData)
	if err != nil {
		panic(err)
	}

	filenames := util.GenerateFilenames(generator.Items, config.FilenameGenFunc)
	var zipData bytes.Buffer
	err = generator.Zip(&zipData, filenames)
	if err != nil {
		panic(err)
	}

	js.CopyBytesToJS(dstJsObject, zipData.Bytes())
	return dstJsObject
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
	fmt.Println("DEBUG: ", jsInputData)
	fmt.Println("DEBUG: ", jsInputData.Type().String())
	dataLen := jsInputData.Length()

	fmt.Println("DEBUG: jsInputData Len:", dataLen)
	inputData := make([]byte, dataLen)
	DEBUG_READ_BYTES := js.CopyBytesToGo(inputData, input.Get("data"))
	fmt.Println("DEBUG: bytes of file copied:", DEBUG_READ_BYTES)
	inputDataReader := bytes.NewReader(inputData)
	items, _, err := util.ParseItems(inputDataReader, csvDelim)
	if err != nil {
		return nil, err
	}
	fmt.Println("DEBUG: Items:", items)

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

func main() {
	js.Global().Set("generate_go", js.FuncOf(Generate))

	<-make(chan bool)
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
