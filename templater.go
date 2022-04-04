package emailtemplater

// TODO: Протестировать gomail, если позволяет формировать черновики, заменить код для
// формирования заголовков и тела на эту библиотеку

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/fedorwk/templater"
	"github.com/pkg/errors"
)

type EmailTemplater struct {
	Subject              string
	Template             []byte
	Items                []map[string]string
	EmailPlaceholder     string
	StartDelim, EndDelim string
}

func (et *EmailTemplater) ParseTemplate(src io.Reader) error {
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
	et.Template = template.Bytes()
	return nil
}

func (et *EmailTemplater) ExecuteToWriters(dests ...io.Writer) error {
	if len(dests) != len(et.Items) {
		return errors.New("number of writers (dests) must be equal to number of Items")
	}
	headers, err := et.buildEMLHeaders()
	if err != nil {
		return errors.Wrap(err, "err while building template:")
	}
	emlTemplate, err := et.buildEMLTemplate(headers)
	if err != nil {
		return errors.Wrap(err, "err while building template:")
	}
	templater := templater.NewTemplater(
		string(emlTemplate),
		templater.NewReplacer(et.Items, et.StartDelim, et.EndDelim),
	)
	for i := 0; i < len(et.Items); i++ {
		_, err = templater.ExecuteToStream(i, dests[i])
		if err != nil {
			return errors.Wrapf(err, "writing %d'th file", i+1)
		}
	}
	return nil
}

func (et *EmailTemplater) buildEMLHeaders() ([]byte, error) {
	if et.EmailPlaceholder == "" {
		et.EmailPlaceholder = DefineEmailPlaceholder(et.Items)
		if et.EmailPlaceholder == "" {
			return nil, ErrNoEmailPlaceholder
		}
	}
	result := bytes.NewBuffer(nil)
	_, err := result.WriteString(fmt.Sprintf(EMLHeaderFormatString, et.EmailPlaceholder, et.Subject))
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

func (et *EmailTemplater) buildEMLTemplate(headers []byte) ([]byte, error) {
	buf := bytes.NewBuffer(headers)
	buf.Grow(len(et.Template) + 1)
	err := buf.WriteByte('\n')
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(et.Template)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
