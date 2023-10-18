package types

type Paginateable[T any] struct {
	CurrentPage uint16 `json:"currentPage"`
	PageSize    uint16 `json:"pageSize"`
	Data        []T    `json:"data"`
}
