package users

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// roles should be list of short roles from usermodel.RoleShort struct
func (jwtM JWTMaker) CreateAccessToken(id int64, userName string, email string, roles []string) (string, *UserClaims, error) {
	return jwtM.createToken(id, userName, email, jwtM.accessTokenDuration, ScopeAuthToken, roles)
}

// roles should be list of short roles from usermodel.RoleShort struct
func (jwtM JWTMaker) createToken(id int64, userName string, email string, duration time.Duration, scope string, roles []string) (string, *UserClaims, error) {
	userClaims, err := jwtM.NewUserClaims(id, userName, email, duration, scope, roles)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenStr, err := token.SignedString([]byte(jwtM.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("%w, email: %s", ErrJWTSignatureFail, email)
	}

	return tokenStr, userClaims, nil
}

// roles should be list of short roles from usermodel.RoleShort struct
func (jwtM JWTMaker) NewUserClaims(id int64, userName string, email string, duration time.Duration, scope string, roles []string) (*UserClaims, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, ErrJWTCreateId
	}
	return &UserClaims{
		Id:         id,
		Email:      email,
		Scope:      scope,
		ShortRoles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Subject:   userName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		}}, nil

}

func (jwtM JWTMaker) VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (any, error) {
		//Verify signing method
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrJWTIncorrectSignMethod
		}
		return []byte(jwtM.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w\n%w", ErrJWTErrorReadingTheToken, err)
	}

	userClaims, ok := token.Claims.(*UserClaims)

	if !ok {
		return nil, ErrJWTInvalidTokenClaims
	}
	return userClaims, nil
}

// Check that access token was passed and not refresh token
func (jwtM JWTMaker) VerifyAccessToken(accessJWT string) (*UserClaims, error) {

	userClaims, err := jwtM.VerifyToken(accessJWT)

	if err != nil {
		return nil, err
	}

	if userClaims.Scope != ScopeAuthToken {
		return nil, ErrNotAnAccessToken
	}

	return userClaims, nil

}

func (jwtM *JWTMaker) VerifyRefreshToken(refreshJWT string) (*UserClaims, error) {
	userClaims, err := jwtM.VerifyToken(refreshJWT)

	if err != nil {
		return nil, err
	}
	if userClaims.Scope != ScopeRefreshToken {
		return nil, ErrNotARefreshToken
	}

	return userClaims, nil
}

func (jwtM JWTMaker) CreateRefreshToken(id int64, userName string, email string, roles []string) (string, *UserClaims, error) {

	return jwtM.createToken(id, userName, email, jwtM.refreshTokenDuration, ScopeRefreshToken, roles)
}
