package users

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/custresponse"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/jsonlog"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/validator"
)

// Reads token from input
// Retrns error and validation problems if any accured
// There are two types of errors:
//  1. ReadJSON error
//  2. TokenInput validation err
func (us *UserService) ReadToken(w http.ResponseWriter, r *http.Request, input *usermodel.TokenInput) (error, map[string]string) {

	err := custresponse.ReadJSON(w, r, input)

	if err != nil {
		return err, nil
	}

	v := validator.New()

	usermodel.ValidateTokenPlainText(v, input.TokenPlainText)

	if !v.Valid() {
		return usermodel.ErrTokenInputValidation, v.Errors
	}

	return nil, nil

}

func (us *UserService) ReadLogin(w http.ResponseWriter, r *http.Request, input *usermodel.UserLoginInput) (error, map[string]string) {
	err := custresponse.ReadJSON(w, r, input)

	if err != nil {
		return err, nil
	}

	v := validator.New()
	usermodel.ValidateUserEmail(v, input.Email)

	if !v.Valid() {
		return usermodel.ErrLoginInputValidation, v.Errors
	}
	return nil, nil
}

func (us *UserService) ActivateRegistrationToken(ctx context.Context, token *usermodel.Token) error {

	isTokenOk := us.checkToken(token)

	if !isTokenOk {
		return fmt.Errorf("ERROR: %w. ", ErrTokenExpired)
	}

	// Deactivate registration token
	err := us.Models.Token.DeactivateTokensForUsers(usermodel.ScopeActivation, token.UserId)
	if err != nil {
		return fmt.Errorf("ERROR: %w. Could not Deactivate Activation token", err)
	}

	err = us.Models.User.ActivateUserById(ctx, token.UserId)

	if err != nil {
		return fmt.Errorf("ERROR: %w. Could not activate user", err)
	}

	return nil
}

func (us *UserService) checkToken(token *usermodel.Token) bool {

	isOk := false

	if time.Now().After(token.Expire) {

		expiredMap := map[string]string{
			"UserId": strconv.FormatInt(token.UserId, 10),
		}
		jsonlog.PrintInfo("Token is expired", expiredMap, nil)

		return isOk

	}

	isOk = true

	return isOk

}

// Returns *User, Token Plain Text, error
func (us *UserService) recreateToken(ctx context.Context, userId int64) (*usermodel.User, string, error) {

	err := us.Models.Token.DeactivateTokensForUsers(usermodel.ScopeActivation, userId)

	if err != nil {
		return nil, "", fmt.Errorf("ERROR: %w. Could not deactivate token", err)
	}

	user, err := us.Models.User.GetById(ctx, userId)

	if err != nil {
		return nil, "", fmt.Errorf("ERROR: %w. Could not get user in Token recreation", err)
	}

	token, err := us.createRegistrationToken(user)

	if err != nil {
		return nil, "", err
	}

	return user, token.PlainText, nil

}
