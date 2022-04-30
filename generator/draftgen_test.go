package generator_test

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/fedorwk/draftgen/generator"
)

func TestParseTemplate(t *testing.T) {
	inputTemplate := bytes.NewBufferString(
		`subject:test
line one
line two`,
	)
	emailteplate := &generator.DraftGenerator{}
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

func TestParseItems(t *testing.T) {
	inputCSV := `good;price
apple;10
orange;20
expensive;3,000`
	gen := generator.DraftGenerator{}
	err := gen.ParseItems(strings.NewReader(inputCSV), ";")
	if err != nil {
		t.Errorf("unexpected error:%s", err)
	}
	want := []map[string]string{
		{"good": "apple", "price": "10"},
		{"good": "orange", "price": "20"},
		{"good": "expensive", "price": "3,000"},
	}
	if !reflect.DeepEqual(gen.Items, want) {
		t.Errorf("unexpected result, \nwant:%v\ngot:%v", want, gen.Items)
	}

	raggedCSV := `good;price
apple;10
orange
expensive;3.000`
	err = gen.ParseItems(strings.NewReader(raggedCSV), ";")
	if err == nil {
		t.Error("ragged input error expected, got nil")
	}
}
