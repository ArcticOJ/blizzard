package types

type Paginateable[T any] struct {
	Count       uint32 `json:"count"`
	CurrentPage uint32 `json:"currentPage"`
	PageSize    uint32 `json:"pageSize"`
	Data        []T    `json:"data"`
}
