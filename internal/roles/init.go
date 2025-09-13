package roles

import (
	"database/sql"

	"github.com/konstantin-suspitsyn/datacomrade/data/rolesmodel"
)

type RoleService struct {
	Models *rolesmodel.Queries
}

func New(db *sql.DB) *RoleService {
	return &RoleService{
		Models: rolesmodel.New(db),
	}
}
