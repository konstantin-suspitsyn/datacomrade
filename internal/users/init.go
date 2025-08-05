package users

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/konstantin-suspitsyn/datacomrade/data/rolesmodel"
	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
)

const ScopeAuthToken = "scope_auth"
const ScopeRefreshToken = "scope_refresh"

type UserService struct {
	UserModels *usermodel.Models
	JWTMaker   *JWTMaker
	RoleModel  *rolesmodel.Queries
}

type JWTMaker struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

type UserClaims struct {
	Id         int64    `json:"id"` // User Id
	Email      string   `json:"email"`
	Scope      string   `json:"scope"`
	ShortRoles []string `json:"roles"` // Short Roles from UserRoles
	jwt.RegisteredClaims
}

// Initializes UserService
func New(db *sql.DB) *UserService {

	models := usermodel.NewModel(db)
	roleModel := rolesmodel.New(db)

	userService := UserService{
		UserModels: &models,
		JWTMaker:   NewJWTMaker(),
		RoleModel:  roleModel,
	}

	return &userService
}

func NewJWTMaker() *JWTMaker {

	jwtConfig := configs.InitJWTConfig()

	return &JWTMaker{
		secretKey:            jwtConfig.SecretKey,
		accessTokenDuration:  jwtConfig.TokenDuration,
		refreshTokenDuration: jwtConfig.RefreshTokenDuration,
	}
}
