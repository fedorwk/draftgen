package emailtemplater

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

type EmailTemplater struct {
	Subject          string
	Template         []byte
	Items            []map[string]string
	EmailPlaceholder string
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

func (et *EmailTemplater) ExecuteToWriters(w ...io.Writer) error {
	headers := et.buildEMLHeaders()
	emlTemplate := et.buildEMLTemplate(headers)
	// templater.NewTemplater -> form eml for all et.Items and write to Writers (w)
	return nil
}

func (et *EmailTemplater) buildEMLHeaders() []byte {
	result := bytes.NewBuffer(nil)
	// take headers from eml template,
	// extend "To:*" part with et.EmailPlaceholder (after check for existance)
	// write to result
	return result.Bytes()
}

func (et *EmailTemplater) buildEMLTemplate(headers []byte) []byte {
	// glue EMLHeaders and Body template (et.Template)
	return nil
}
