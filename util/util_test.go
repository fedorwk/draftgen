package util_test

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/fedorwk/draftgen/util"
)

func ExampleGenerateFilenames() {
	items := []map[string]string{
		{"good": "apple", "price": "10"},
		{"good": "orange", "price": "20"},
	}
	suffix := ".item"

	// filename will be generated using the following pattern:
	// ["good" field value][index+1].item
	fn := func(index int, item map[string]string) string {
		filename := item["good"] + strconv.Itoa(index+1) + suffix
		return filename
	}

	filenames := util.GenerateFilenames(items, fn)
	fmt.Printf("%+v\n", filenames)
	// Output: [apple1.item orange2.item]
}

func TestParseItems(t *testing.T) {
	inputCSV := `good;price
apple;10
orange;20
expensive;3,000`
	got, _, err := util.ParseItems(strings.NewReader(inputCSV), ";")
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
	_, _, err = util.ParseItems(strings.NewReader(raggedCSV), ";")
	if err == nil {
		t.Error("ragged input error expected, got nil")
	}
}
