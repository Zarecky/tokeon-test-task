package customerror

type BadRequestError struct {
	Message string
}

func (r BadRequestError) Error() string {
	return r.Message
}

type ForbiddenError struct {
	Message string
}

func (r ForbiddenError) Error() string {
	return r.Message
}
