package jikan

type PaginationItems struct {
	Count   int `json:"count"`
	Total   int `json:"total"`
	PerPage int `json:"per_page"`
}

type Pagination struct {
	LastVisiblePage int             `json:"last_visible_page"`
	HasNextPage     bool            `json:"has_next_page"`
	CurrentPage     int             `json:"current_page"`
	Items           PaginationItems `json:"items"`
}

type Result[T any] struct {
	Data T `json:"data"`
}

type PaginatedResults[T any] struct {
	Data       T `json:"data"`
	Pagination `json:"pagination"`
}
