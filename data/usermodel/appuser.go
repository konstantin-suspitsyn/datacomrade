package usermodel

type AppUser struct {
	Id         int64
	UserName   string
	ShortRoles []string
}

var AnonymousUser = &AppUser{}

func (au *AppUser) IsAnonimous() bool {
	return AnonymousUser == au
}
