package shareddata

import (
	"context"
	"net/http"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/konstantin-suspitsyn/datacomrade/data/paginationmodel"
	"github.com/konstantin-suspitsyn/datacomrade/data/sharedmodels"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/urlparams"
)

func (sds *SharedDataService) getDataWithPager(r *http.Request) (*sharedmodels.DomainsWithPagerDTO, error) {
	ctx := r.Context()

	urlPagingParams, err := urlparams.GetPager(r)
	if err != nil {
		return nil, err
	}

	paginator, err := sds.generatePaginator(ctx, urlPagingParams)

	if err != nil {
		return nil, err
	}
	argsForPager := sharedmodels.GetDomainsWithPagerParams{
		Limit:  paginator.GetLimit(),
		Offset: paginator.GetOffset(),
	}
	rows, err := sds.Models.GetDomainsWithPager(ctx, argsForPager)

	if err != nil {
		return nil, err
	}

	dto := sharedmodels.DomainsWithPagerDTO{
		Data:   rows,
		Paging: paginator,
	}

	return &dto, nil
}

func (sds *SharedDataService) generatePaginator(ctx context.Context, pager *urlparams.Pager) (*paginationmodel.Pagination, error) {

	totalItems, err := sds.Models.CountActiveRows(ctx)

	if err != nil {
		return nil, err
	}

	paginator, err := paginationmodel.New(totalItems, pager.PageSize, pager.CurrentPage, configs.DOMAIN_LINK)

	if err != nil {
		return nil, err
	}

	return paginator, nil
}
