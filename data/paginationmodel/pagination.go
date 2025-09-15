package paginationmodel

import (
	"fmt"
	"math"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
)

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

func New(totalItems int64, pageSize int, currentPage int64, link string) (*Pagination, error) {
	// if pageSize == 0, then use standart pager and do not use it in link
	pageSizeUpd := configs.ITEMS_PER_PAGE

	if pageSize != 0 {
		pageSizeUpd = pageSize
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(pageSizeUpd)))

	first := ""
	previous := ""
	next := ""
	last := ""

	if currentPage < totalPages {
		if currentPage != 1 {
			first = fmt.Sprintf("%s?%s=%d", link, configs.PAGE_PARAM, 1)
			previous = fmt.Sprintf("%s?%s=%d", link, configs.PAGE_PARAM, currentPage-1)
		}
		next = fmt.Sprintf("%s?%s=%d", link, configs.PAGE_PARAM, currentPage+1)
		last = fmt.Sprintf("%s?%s=%d", link, configs.PAGE_PARAM, totalPages)
	}

	if currentPage > totalPages {
		return nil, fmt.Errorf("ERROR. %w. Current page is %d. Total pages are %d", ErrCurrentPageIsBiggerThanTotal, currentPage, totalPages)
	}

	// Adding pageSize
	first, err := addPagesiseToLink(first, pageSize)
	if err != nil {
		return nil, fmt.Errorf("%w. First page", ErrNegativaPageSize)
	}
	previous, err = addPagesiseToLink(previous, pageSize)
	if err != nil {
		return nil, fmt.Errorf("%w. Previous page", ErrNegativaPageSize)
	}
	next, err = addPagesiseToLink(next, pageSize)
	if err != nil {
		return nil, fmt.Errorf("%w. Next page", ErrNegativaPageSize)
	}
	last, err = addPagesiseToLink(last, pageSize)
	if err != nil {
		return nil, fmt.Errorf("%w. Last page", ErrNegativaPageSize)
	}

	return &Pagination{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    pageSizeUpd,
		First:       first,
		Previous:    previous,
		Next:        next,
		Last:        last,
	}, nil
}

func addPagesiseToLink(link string, pageSize int) (string, error) {

	// Separator if link == ""
	// To add link style of link?sort=
	separator := "?"

	if pageSize < 0 {
		return "", fmt.Errorf("%w", ErrNegativaPageSize)
	}

	if pageSize == 0 {
		return link, nil
	}

	if link != "" {
		separator = "&"
	}

	return fmt.Sprintf("%s%s%s=%d", link, separator, configs.ITEMS_PER_PAGE_PARAM, pageSize), nil
}

func (p *Pagination) GetOffset() int64 {
	return (p.CurrentPage - 1) * int64(p.PageSize)
}

func (p *Pagination) GetLimit() int64 {
	return int64(p.PageSize)
}
