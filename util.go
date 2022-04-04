package emailtemplater

import "strings"

var EMLHeaderFormatString = `MIME-Version: 1.0
X-Unsent: 1
To: %s
Subject: %s
Content-Type: multipart/alternative
`

func DefineEmailPlaceholder(items []map[string]string) string {
	for key := range items[0] {
		switch strings.ToLower(key) {
		case "email", "e-mail", "mail":
			return key
		}
	}
	return ""
}
