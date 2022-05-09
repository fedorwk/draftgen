package server

import (
	"bytes"
	"log"
	"net/http"

	"github.com/fedorwk/draftgen/generator"
	"github.com/fedorwk/draftgen/util"
	"github.com/fedorwk/go-util/data/delimiterdetector"
)

func serviceHTMLHandlerFn(w http.ResponseWriter, r *http.Request) {
	if EmbeddedHTML == nil {
		log.Println(ErrNilHTML)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(EmbeddedHTML)
}

func generateDraftsHandlerFn(w http.ResponseWriter, r *http.Request) {
	gen, err := generatorFromRequest(r)
	if err != nil {
		log.Println(err)
		return
	}
	filenames := util.GenerateFilenames(gen.Items, config.App.FilenameGenFunc)

	w.Header().Set("Content-Disposition", "attachment; filename="+config.App.OutputFilename)
	err = gen.Zip(w, filenames)
	if err != nil {
		log.Println(err)
	}
}

func generatorFromRequest(r *http.Request) (*generator.DraftGenerator, error) {
	err := r.ParseMultipartForm(10 << 20) // max request size of 10 MB
	if err != nil {
		return nil, err
	}
	items, err := itemsFromRequest(r)
	if err != nil {
		return nil, err
	}
	subject := r.Form.Get(config.HTTPForm.Subject)
	template := r.Form.Get(config.HTTPForm.Template)
	emailPlaceholder := r.Form.Get(config.HTTPForm.EmailPlaceholder)
	startDelim := r.Form.Get(config.HTTPForm.StartDelim)
	endDelim := r.Form.Get(config.HTTPForm.EndDelim)

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

func itemsFromRequest(r *http.Request) ([]map[string]string, error) {
	file, _, err := r.FormFile(config.HTTPForm.Data)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var raw bytes.Buffer
	_, err = raw.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	csvDelim := r.Form.Get(config.HTTPForm.CSVDelim)
	if csvDelim == "" {
		csvDelim, err = delimiterdetector.Parse(
			bytes.NewReader(raw.Bytes()),
			config.App.LinesCountToAnalyzeCSV,
		)
		if err != nil {
			return nil, err
		}
	}

	items, _, err := util.ParseItems(&raw, csvDelim)
	if err != nil {
		return nil, err
	}
	return items, nil
}
