package emailtemplater

import "strings"

func DefineEmailPlaceholder(items []map[string]string) string {
	for key := range items[0] {
		switch strings.ToLower(key) {
		case "email", "e-mail", "mail":
			return key
		}
	}
	return ""
}
