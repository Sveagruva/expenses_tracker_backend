package repository

type Pagination struct {
	Page  int64 `json:"page" binding:"required"`
	Items int64 `json:"items" binding:"required"`
}

type PaginationResponse[T any] struct {
	Items []T   `json:"items"`
	Count int64 `json:"count"`
}

type SqlPagination struct {
	Limit  int64
	Offset int64
}

func ResolvePagination(pagination *Pagination) SqlPagination {
	return SqlPagination{
		Limit:  pagination.Items,
		Offset: (pagination.Page - 1) * pagination.Items,
	}
}
