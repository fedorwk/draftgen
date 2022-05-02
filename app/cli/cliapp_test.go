package cli

import (
	"os"
	"testing"
)

func TestRunFiles(t *testing.T) {
	os.Args = append(os.Args,
		"-t", "testdata/testTemplate",
		"-d", "testdata/testCSV")
	err := Run()
	if err != nil {
		t.Error(err)
	}
}

func TestRunZip(t *testing.T) {
	os.Args = append(os.Args,
		"-t", "testdata/testTemplate",
		"-d", "testdata/testCSV",
		"-z")
	err := Run()
	if err != nil {
		t.Error(err)
	}
}
