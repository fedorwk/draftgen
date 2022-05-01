package util_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/fedorwk/draftgen/util"
)

func TestParseItems(t *testing.T) {
	inputCSV := `good;price
apple;10
orange;20
expensive;3,000`
	got, err := util.ParseItems(strings.NewReader(inputCSV), ";")
	if err != nil {
		t.Errorf("unexpected error:%s", err)
	}
	want := []map[string]string{
		{"good": "apple", "price": "10"},
		{"good": "orange", "price": "20"},
		{"good": "expensive", "price": "3,000"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected result, \nwant:%v\ngot:%v", want, got)
	}

	raggedCSV := `good;price
apple;10
orange
expensive;3.000`
	_, err = util.ParseItems(strings.NewReader(raggedCSV), ";")
	if err == nil {
		t.Error("ragged input error expected, got nil")
	}
}
