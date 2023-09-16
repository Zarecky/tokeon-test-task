package common

import "errors"

var (
	ErrNotFound  = errors.New("not found")
	ErrEmptyID   = errors.New("empty id")
	ErrEmptyCode = errors.New("empty code")
)
