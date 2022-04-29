package util

import (
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
