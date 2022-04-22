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

func IdentifyCSVDelimiter() {
	// TODO:
	// создать хэшмап из стандартных разделителей:инт
	// посчитать количество вхождений каждого ключа в первую сторку
	// если колво вхождений ключа в строку - 0 - удалить этот ключ
	// если колво вхождений всех ключей 0 - разделителей нет или нестандартный
	// для каждой следующей строки:
	// если вхождений ключа меньше - исключить ключ из мапы
	// если остался один ключ - проверить, чтобы в следующей строке совпадало количество
	// вхождений.
	// если количество вхождений единственного ключа во всех строках совпадает -
	// это и есть ключ
}
