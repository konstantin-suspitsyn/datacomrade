package rolesmodel

import (
	"database/sql"
	"time"
)

type Action struct {
	Id          int64
	Name        string
	Description string
	IsDeleted   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Resource struct {
	Id          int64
	Name        string
	Description string
	IsDeleted   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Role struct {
	Id            int64
	RoleNameLong  string
	RoleNameShort string
	Description   string
	JwtExport     bool
	IsDeleted     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserRole struct {
	Id        int64
	UserId    int64
	RoleId    int64
	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserAccess struct {
	Id             int64
	UserId         int64
	ResourceId     int64
	ResourceTypeId int64
	ActionId       int64
	IsDeleted      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RoleAccess struct {
	Id             int64
	RoleId         int64
	ResourceId     int64
	ResourceTypeId int64
	ActionId       int64
	IsDeleted      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RoleModel struct {
	DB *sql.DB
}
