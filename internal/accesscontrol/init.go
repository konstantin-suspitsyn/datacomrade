// This is service to setup domaind, roles and other auth staff
package accesscontrol

import (
	"database/sql"

	"github.com/konstantin-suspitsyn/datacomrade/internal/shareddata"
)

type AccessControlService struct {
	SharedDataService *shareddata.SharedDataService
}

func New(db *sql.DB) *AccessControlService {
	sharedDataService := shareddata.New(db)
	return &AccessControlService{
		SharedDataService: sharedDataService,
	}
}
