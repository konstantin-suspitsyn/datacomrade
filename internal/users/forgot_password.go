package users

import (
	"context"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/validator"
)

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	// Validate email
	v := validator.New()

	usermodel.ValidateUserEmail(v, email)
	if !v.Valid() {
		return nil, fmt.Errorf("%w. Email: %s", usermodel.ErrEmailValidation, email)
	}
	// Get User
	return us.Models.User.GetByEmail(ctx, email)
}

func (us *UserService) CreateForgotPasswordToken(user *usermodel.User) (*usermodel.Token, error) {

	token, err := us.Models.Token.New(user.Id, configs.TokenDuration, usermodel.ScopeForgotPassword)
	if err != nil {
		return nil, err
	}

	return token, nil
}
