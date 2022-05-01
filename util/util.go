package util

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

func DefineEmailPlaceholder(items []map[string]string) string {
	for key := range items[0] {
		switch strings.ToLower(key) {
		case "email", "e-mail", "mail":
			return key
		}
	}
	return ""
}

type NameGenFn func(index int, item map[string]string) string

func GenerateFilenames(items []map[string]string, fn NameGenFn) []string {
	filenames := make([]string, len(items))
	for index, item := range items {
		filenames[index] = fn(index, item)
	}
	return filenames
}

func ParseItems(csv io.Reader, delimiter string) ([]map[string]string, error) {
	scanner := bufio.NewScanner(csv)
	var headers []string
	if ok := scanner.Scan(); ok {
		headers = strings.Split(scanner.Text(), delimiter)
	}

	items := make([]map[string]string, 0)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), delimiter)
		if len(values) != len(headers) {
			return nil, errors.New("ragged csv input")
		}
		item := make(map[string]string, len(headers))
		for i, header := range headers {
			item[header] = values[i]
		}
		items = append(items, item)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return items, nil
}
