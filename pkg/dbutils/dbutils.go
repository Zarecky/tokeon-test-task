package dbutils

import "fmt"

type Option func(*Options)

type Options struct {
	MaxLimit uint64
}

func WithMaxLimit(v uint64) Option {
	return func(o *Options) {
		o.MaxLimit = v
	}
}

func Pagination(page, pageSize *uint64, opts ...Option) (limit, offset uint64, err error) {
	options := Options{
		MaxLimit: 100,
	}

	for _, opt := range opts {
		opt(&options)
	}

	limit = 100
	offset = 0

	if pageSize != nil {
		limit = *pageSize
	}

	if page != nil && *page != 0 {
		offset = (*page - 1) * limit
	}

	if limit < 1 {
		return 0, 0, fmt.Errorf("pagination error: page size cannot be less than 1")
	}

	if limit > options.MaxLimit {
		return 0, 0, fmt.Errorf("pagination error: page size cannot be greater than %d", options.MaxLimit)
	}

	return
}
