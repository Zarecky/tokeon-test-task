package common

type PageOptions struct {
	Page     *uint64
	PageSize *uint64
}

type FindResponse[T any] struct {
	Items []T `json:"items"`
}

type FindResponseWithCount[T any] struct {
	Items []T   `json:"items"`
	Count int64 `json:"count"`
}

func GetLikeVal(v string) string {
	return "%" + v + "%"
}

type CommonFindRequest struct {
	OrderBy            string  `query:"order_by" json:"order_by"`
	QueryIsAscOrdering bool    `query:"is_asc_ordering" json:"is_asc_ordering"`
	Page               *uint64 `query:"page" json:"page"`
	PageSize           *uint64 `query:"page_size" son:"page_size"`
}
