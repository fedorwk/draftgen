package draftgen_test

import (
	"bytes"
	"testing"

	"github.com/fedorwk/draftgen"
)

func TestParseTemplate(t *testing.T) {
	inputTemplate := bytes.NewBufferString(
		`subject:test
line one
line two`,
	)
	emailteplate := &draftgen.DraftGenerator{}
	err := emailteplate.ParseTemplate(inputTemplate)
	if err != nil {
		t.Error(err)
	}
	if emailteplate.Subject != "test" {
		t.Errorf(`wrong subject parsed, want: "text", got: "%s"`, emailteplate.Subject)
	}
	if string(emailteplate.Template) != "line one\nline two\n" {
		t.Errorf("wrong template parsed, want:\nline one\nline two\ngot:\n%s", emailteplate.Template)
	}
}
