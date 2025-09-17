package shareddata

import (
	"context"
	"database/sql"
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

func (sds *SharedDataService) CreateDomain(ctx context.Context, name string, description string, userId int64) (sharedmodels.SharedDomain, error) {
	descriptionForDomain := sql.NullString{
		String: description,
		Valid:  true,
	}

	createDomainParams := sharedmodels.CreateDomainParams{
		Name:        name,
		Description: descriptionForDomain,
		UserID:      userId,
	}

	shareDomain, err := sds.Models.CreateDomain(ctx, createDomainParams)

	return shareDomain, err
}

func (sds *SharedDataService) UpdateDomainByUser(ctx context.Context, domainId int64, name string, description string, userId int64) (sharedmodels.SharedDomain, error) {
	descriptionForDomain := sql.NullString{
		String: description,
		Valid:  true,
	}

	updateOneParams := sharedmodels.UpdateOneParams{
		Name:        name,
		Description: descriptionForDomain,
		UserID:      userId,
		ID:          domainId,
	}

	sd, err := sds.Models.UpdateOne(ctx, updateOneParams)
	return sd, err
}

func (sds *SharedDataService) DeleteDomain(ctx context.Context, domainId int64, userId int64) error {
	deleteDomainParams := sharedmodels.DeleteDomainParams{
		ID:     domainId,
		UserID: userId,
	}

	return sds.Models.DeleteDomain(ctx, deleteDomainParams)
}
