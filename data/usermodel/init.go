package usermodel

import (
	"database/sql"
)

type Models struct {
	Token        TokenModel
	User         UserModel
	RefreshToken RefreshTokenModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Token:        TokenModel{DB: db},
		User:         UserModel{DB: db},
		RefreshToken: RefreshTokenModel{DB: db},
	}
}
