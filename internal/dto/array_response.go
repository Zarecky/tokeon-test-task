package dto

type ArrayResponse[T any] struct {
	Items []T `json:"items"`
}

type ArrayWithAmountResponse[T any] struct {
	Items []T   `json:"items"`
	Count int64 `json:"count"`
}
