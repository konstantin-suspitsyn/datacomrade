package shareddata

import (
	"database/sql"

	"github.com/konstantin-suspitsyn/datacomrade/data/sharedmodels"
)

type SharedDataService struct {
	Models *sharedmodels.Queries
}

func New(db *sql.DB) *SharedDataService {
	return &SharedDataService{
		Models: sharedmodels.New(db),
	}

}
