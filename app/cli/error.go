package cli

import "errors"

var (
	ErrNoTemplateArg = errors.New("no template file passed")
	ErrNoDataArg     = errors.New("no data source passed")
)
