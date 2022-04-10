package draftgen

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	email "github.com/fedorwk/jw-email"
	"github.com/fedorwk/templater"
	"github.com/pkg/errors"
)

type DraftGenerator struct {
	Subject              string
	Template             string
	Items                []map[string]string
	EmailPlaceholder     string
	StartDelim, EndDelim string

	cachedTemplater *templater.Templater
}

func (et *DraftGenerator) ParseTemplate(src io.Reader) error {
	scanner := bufio.NewScanner(src)
	ok := scanner.Scan()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	if !ok {
		return errors.New("empty source")
	}
	template := bytes.NewBuffer(nil)
	if strings.HasPrefix(scanner.Text(), "subject:") {
		et.Subject = strings.TrimPrefix(scanner.Text(), "subject:")
	} else {
		template.Write(scanner.Bytes())
	}
	for scanner.Scan() {
		template.Write(scanner.Bytes())
		template.WriteByte('\n')
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	et.Template = template.String()
	return nil
}

func (et *DraftGenerator) Execute(itemIndex int, dst io.Writer) error {
	// if there was no cached Templater it is created and deleted at the end of the function execution
	// used while singel Execute call
	if et.cachedTemplater == nil {
		et.cachedTemplater = templater.NewTemplater(et.Template, et.Items, et.StartDelim, et.EndDelim)
		defer func() { et.cachedTemplater = nil }()
	}
	item := et.Items[itemIndex]

	eml := email.NewEmail()
	eml.Draft = true
	eml.To = []string{item[et.EmailPlaceholder]}
	eml.Subject = et.Subject
	eml.HTML = []byte(et.cachedTemplater.ExecuteToString(itemIndex))

	rawbytes, err := eml.Bytes()
	if err != nil {
		return err
	}
	_, err = dst.Write(rawbytes)
	if err != nil {
		return err
	}
	return nil
}

func (et *DraftGenerator) ExecuteAll(dests ...io.Writer) error {
	if len(dests) != len(et.Items) {
		return errors.New("number of writers (dests) must be equal to number of Items in receiver")
	}
	if et.EmailPlaceholder == "" {
		err := et.defineEmailPlaceholder()
		if err != nil {
			return err
		}
	}
	et.cachedTemplater = templater.NewTemplater(et.Template, et.Items, et.StartDelim, et.EndDelim)
	defer func() { et.cachedTemplater = nil }()

	for i := range et.Items {
		err := et.Execute(i, dests[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (et *DraftGenerator) defineEmailPlaceholder() error {
	if placeHolder := DefineEmailPlaceholder(et.Items); placeHolder == "" {
		return ErrNoEmailPlaceholder
	} else {
		et.EmailPlaceholder = placeHolder
	}
	return nil
}
