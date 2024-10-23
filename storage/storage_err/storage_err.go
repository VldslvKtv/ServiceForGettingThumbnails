package storage_err

import "errors"

var (
	ErrNotFoundCache = errors.New("no cache")
)
