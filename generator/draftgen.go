package generator

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/fedorwk/draftgen/util"
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

func (dg *DraftGenerator) ParseTemplate(src io.Reader) error {
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
		if dg.Subject == "" {
			dg.Subject = strings.TrimPrefix(scanner.Text(), "subject:")
		}
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
	dg.Template = template.String()
	return nil
}

func (dg *DraftGenerator) Execute(itemIndex int, dst io.Writer) error {
	// if there was no cached Templater it is created and deleted at the end of the function execution
	// used while singel Execute call
	if dg.cachedTemplater == nil {
		dg.cachedTemplater = templater.NewTemplater(dg.Template, dg.Items, dg.StartDelim, dg.EndDelim)
	}
	if dg.EmailPlaceholder == "" {
		err := dg.defineEmailPlaceholder()
		if err != nil {
			return err
		}
	}

	item := dg.Items[itemIndex]

	eml := email.NewEmail()
	eml.Draft = true
	eml.To = []string{item[dg.EmailPlaceholder]}
	eml.Subject = dg.Subject
	eml.HTML = []byte(dg.cachedTemplater.ExecuteToString(itemIndex))

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

func (dg *DraftGenerator) ExecuteAll(dsts ...io.Writer) error {
	if len(dsts) != len(dg.Items) {
		return errors.New("number of writers (dests) must be equal to number of Items in receiver")
	}
	dg.cachedTemplater = templater.NewTemplater(dg.Template, dg.Items, dg.StartDelim, dg.EndDelim)
	defer func() { dg.cachedTemplater = nil }()

	for i := range dg.Items {
		err := dg.Execute(i, dsts[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (dg *DraftGenerator) defineEmailPlaceholder() error {
	if dg.Items == nil {
		return ErrNoItems
	}
	if placeHolder := util.DefineEmailPlaceholder(dg.Items); placeHolder == "" {
		return ErrNoEmailPlaceholder
	} else {
		dg.EmailPlaceholder = placeHolder
	}
	return nil
}
