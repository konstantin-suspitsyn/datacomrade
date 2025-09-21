package sharedmodels

import (
	"context"
)

type SharedModelsInterface interface {
	CountActiveRows(context.Context) (int64, error)
	CreateDomain(context.Context, CreateDomainParams) (SharedDomain, error)
	DeleteDomain(context.Context, DeleteDomainParams) error
	GetDomainById(context.Context, int64) (GetDomainByIdRow, error)
	GetDomains(context.Context) ([]GetDomainsRow, error)
	GetDomainsWithPager(context.Context, GetDomainsWithPagerParams) ([]GetDomainsWithPagerRow, error)
	UpdateOne(context.Context, UpdateOneParams) (SharedDomain, error)
}
