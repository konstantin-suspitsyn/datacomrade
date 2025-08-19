package users

import (
	"context"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/data/rolesmodel"
	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
)

// Checks email and password against database and returns true if user exists
// returns false and error if something wrong with user
func (us *UserService) checkEmailPasswordAndActive(ctx context.Context, email string, password string) (bool, error) {

	user, err := us.UserModels.User.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	// Check if user is activated
	if !user.Activated {
		return false, ErrUserNotActivated
	}

	passwordMatch, err := user.Password.Matches(password)

	if err != nil {
		return false, err
	}

	return passwordMatch, nil
}

// Returns user model, list of user roles, error
func (us *UserService) findUserAndJWTRoles(ctx context.Context, email string) (*usermodel.User, []rolesmodel.GetJWTShortRolesByUserIdRow, error) {

	user, err := us.UserModels.User.GetActiveByEmail(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting user by email %s. ERROR: %w", email, err)
	}

	roles, err := us.RoleModel.GetJWTShortRolesByUserId(ctx, user.Id)

	if err != nil {
		return nil, nil, fmt.Errorf("Error getting roles for user %d. ERROR: %w", user.Id, err)
	}

	return user, roles, nil
}

func (us *UserService) accessAndRefreshTokens(ctx context.Context, user *usermodel.User, roles []rolesmodel.GetJWTShortRolesByUserIdRow) (*usermodel.LoginDTO, error) {

	var rolesString []string

	for _, role := range roles {
		rolesString = append(rolesString, role.RoleNameShort)
	}

	accessToken, userClaims, err := us.JWTMaker.CreateAccessToken(user.Id, user.Name, user.Email, rolesString)

	if err != nil {
		return nil, err
	}

	refreshTokenString, refreshClaims, err := us.JWTMaker.CreateRefreshToken(user.Id, user.Name, user.Email, rolesString)

	if err != nil {
		return nil, err
	}

	refreshToken := usermodel.RefreshToken{
		Id:           refreshClaims.RegisteredClaims.ID,
		UserId:       user.Id,
		Expire:       refreshClaims.RegisteredClaims.ExpiresAt.Time,
		CreatedAt:    refreshClaims.RegisteredClaims.IssuedAt.Time,
		IsActive:     true,
		RefreshToken: refreshTokenString,
	}

	err = us.UserModels.RefreshToken.Insert(ctx, &refreshToken)

	if err != nil {
		return nil, fmt.Errorf("Refresh token creation failed. %w", err)
	}

	loginDTO := usermodel.LoginDTO{
		AccessToken:                accessToken,
		AccessTokenExpirationTime:  userClaims.ExpiresAt.Time,
		RefreshToken:               refreshTokenString,
		RefreshTokenExpirationTime: refreshToken.Expire,
		SessionId:                  refreshToken.Id,
	}

	return &loginDTO, nil
}
