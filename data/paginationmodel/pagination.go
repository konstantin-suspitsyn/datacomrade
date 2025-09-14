package paginationmodel

type Pagination struct {
	TotalItems  int64  `json:"total_items"`
	TotalPages  int64  `json:"total_pages"`
	CurrentPage int64  `json:"current_page"`
	PageSize    int    `json:"page_size"`
	First       string `json:"first"`
	Previous    string `json:"previous"`
	Next        string `json:"next"`
	Last        string `json:"last"`
}

func New(totalItems int64, pageSize int, currentPage int64) *Pagination {
	// TODO: Add First, Previous, Next, Last
	totalPages := totalItems / currentPage
	return &Pagination{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    pageSize,
	}
}
