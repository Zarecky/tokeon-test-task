package errors

import (
	e "errors"
	"fmt"
)

type ValidationError error

func NewValidationError(message string) ValidationError {
	return ValidationError(fmt.Errorf(message))
}

var ErrUserIsNil = e.New("user is nil")
var ErrNotFound = e.New("not found")
var ErrInvalidParams = e.New("invalid params")
var ErrMissingRequiredFields = e.New("missing required fields")
var ErrNotImplemented = e.New("not implemented")
var ErrAlreadyExists = e.New("already exists")
