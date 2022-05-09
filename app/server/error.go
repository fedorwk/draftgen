package server

import "errors"

var (
	ErrNilHTML       = errors.New("EmbeddedHTML not initialized")
	ErrNoItemsParsed = errors.New("no items parsed from data")
)
