package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/validator"
)

// Creates user from payload and tries to validate it
// returns User model, error, map of validation errors
func (us *UserService) createUser(input usermodel.UserRegisterInput) (usermodel.User, error, map[string]string) {
	user := usermodel.User{
		Name:      input.Name,
		Email:     strings.ToLower(input.Email),
		Activated: false,
	}

	user.Password.Set(input.Password)
	v := validator.New()

	if usermodel.ValidateUser(v, &user); !v.Valid() {
		return usermodel.User{}, fmt.Errorf("ERROR: %w. Checking model: User. Input: {Name: %s, Email: %s}", validator.ErrValidation, input.Name, input.Email), v.Errors
	}

	return user, nil, nil

}

func (us *UserService) insertUserToDB(ctx context.Context, user *usermodel.User) (error, map[string]string) {
	err := us.UserModels.User.Insert(ctx, user)

	if err != nil {

		switch {

		case errors.Is(err, usermodel.ErrDuplicateEmail):
			errorMap := map[string]string{
				"email": "Email Already exists",
			}
			return fmt.Errorf("ERROR: %w. Email: %s", usermodel.ErrDuplicateEmail, user.Email), errorMap
		case errors.Is(err, usermodel.ErrDuplicateName):
			errorMap := map[string]string{
				"name": "Name Already exists",
			}
			return fmt.Errorf("ERROR: %w. Name: %s", usermodel.ErrDuplicateEmail, user.Name), errorMap

		default:
			return err, nil
		}
	}

	return nil, nil

}

func (us *UserService) createRegistrationToken(user *usermodel.User) (*usermodel.Token, error) {

	token, err := us.UserModels.Token.New(user.Id, configs.TokenDuration, usermodel.ScopeActivation)
	return token, err

}
