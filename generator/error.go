package generator

import "errors"

var (
	ErrNoEmailPlaceholder = errors.New("no email placeholder provided")
	ErrNoItems            = errors.New("no items passed to generator")
)
