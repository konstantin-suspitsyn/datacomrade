package sharedmodels

import "github.com/konstantin-suspitsyn/datacomrade/data/paginationmodel"

type DomainsWithPagerDTO struct {
	Data   []GetDomainsWithPagerRow    `json:"data"`
	Paging *paginationmodel.Pagination `json:"paging"`
}
