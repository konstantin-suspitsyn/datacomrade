package rolesmodel

import "github.com/konstantin-suspitsyn/datacomrade/data/paginationmodel"

type RolesWithPager struct {
	Data   []GetRolesWithPagerRow      `json:"data"`
	Paging *paginationmodel.Pagination `json:"paging"`
}

type RoleInputDTO struct {
	NameLong    string `json:"name_long"`
	NameShort   string `json:"name_short"`
	Description string `json:"description"`
	JWTExport   bool   `json:"jwt_export"`
}
