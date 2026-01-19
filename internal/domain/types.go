package domain

type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

type SortOptions struct {
	Field string
	Order SortOrder
}

type QueryOptions struct {
	Limit  uint64
	Offset uint64
	Sort   *SortOptions
}
