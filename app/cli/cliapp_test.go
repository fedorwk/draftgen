package cli

import (
	"testing"
)

func TestRunFiles(t *testing.T) {
	cliArgs := []string{
		"-t", "testdata/testTemplate",
		"-d", "testdata/testCSV",
	}
	err := Run(cliArgs)
	if err != nil {
		t.Error(err)
	}
}

func TestRunZip(t *testing.T) {
	cliArgs := []string{
		"-t", "testdata/testTemplate",
		"-d", "testdata/testCSV",
		"-z",
	}
	err := Run(cliArgs)
	if err != nil {
		t.Error(err)
	}
}
