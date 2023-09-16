package validation

import "github.com/go-playground/validator/v10"

type ValidationError struct {
	Err error
}

func (r ValidationError) Error() string {
	return r.Err.Error()
}

type Validator struct {
	validator.Validate
}

func New() *Validator {
	validator := validator.New()
	return &Validator{
		Validate: *validator,
	}
}

func (v *Validator) Struct(s interface{}) error {
	return ValidationError{Err: v.Validate.Struct(s)}
}
