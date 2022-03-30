package emailtemplater

type Email struct {
	To      []string
	CC      []string
	BCC     []string
	Subject string
	Body    []byte
}
