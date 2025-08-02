package users

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	_ "github.com/go-chi/chi/v5"
	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/custresponse"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/shared"
)

func (us *UserService) UserRegister(w http.ResponseWriter, r *http.Request) {

	// Read data from input payload
	var input usermodel.UserRegisterInput

	err := custresponse.ReadJSON(w, r, &input)

	if err != nil {
		custresponse.BadRequestResponse(w, r, fmt.Errorf("ERROR: ReadJson UserRegisterInput. %w", err))
	}

	// Create user
	user, err, errMap := us.createUser(input)

	if err != nil {
		custresponse.FailedValidationResponse(w, r, err, errMap)
		return
	}

	ctx := r.Context()

	err, errMap = us.insertUserToDB(ctx, &user)

	if err != nil {
		switch {
		case errors.Is(err, usermodel.ErrDuplicateEmail):
			custresponse.FailedValidationResponse(w, r, err, errMap)
			return
		case errors.Is(err, usermodel.ErrDuplicateName):
			custresponse.FailedValidationResponse(w, r, err, errMap)
			return

		default:
			custresponse.ServerErrorResponse(w, r, err)
			return
		}
	}

	// Create Registration Token
	token, err := us.createRegistrationToken(&user)

	if err != nil {
		custresponse.ServerErrorResponse(w, r, err)
		return
	}

	// Send email in background
	var wg sync.WaitGroup

	shared.BackroundJob(func() {

		us.SendActivationToken(token.PlainText, user.Email)
	}, &wg)

	userResponceMessage := make(map[string]usermodel.User)
	userResponceMessage["user"] = user

	// Response Token
	err = custresponse.WriteJSON(w, http.StatusAccepted, userResponceMessage, nil)
	if err != nil {
		custresponse.ServerErrorResponse(w, r, err)
	}

	wg.Wait()
}

func (us *UserService) UserActivate(w http.ResponseWriter, r *http.Request) {

	var input usermodel.TokenInput

	err, validationProblems := us.ReadToken(w, r, &input)
	if err != nil {
		switch {
		case errors.Is(err, usermodel.ErrTokenInputValidation):
			custresponse.FailedValidationResponse(w, r, err, validationProblems)
			return
		default:
			custresponse.ServerErrorResponse(w, r, err)
			return
		}
	}

	token, err := us.Models.Token.GetByPlainText(input.TokenPlainText, usermodel.ScopeActivation)

	if err != nil {
		custresponse.ServerErrorResponse(w, r, fmt.Errorf("ERROR: %w. Getting Token from input", err))
		return
	}
	ctx := r.Context()
	err = us.ActivateRegistrationToken(ctx, token)

	if err != nil {
		// TODO: Change types of errors
		switch {
		// Did not find Token
		case errors.Is(err, usermodel.ErrTokenRecordNotFound):
			custresponse.ServerErrorResponse(w, r, err)
			return
		// Token is expired
		case errors.Is(err, ErrTokenExpired):
			custresponse.ServerErrorResponse(w, r, err)
			userUpd, tokenPlainText, recreateTokenErr := us.recreateToken(ctx, token.UserId)
			if recreateTokenErr != nil {
				custresponse.ServerErrorResponse(w, r, recreateTokenErr)
				return
			}
			// Send email in background
			var wg sync.WaitGroup

			shared.BackroundJob(func() {
				us.SendActivationToken(tokenPlainText, userUpd.Email)
			}, &wg)

			// Response Token
			err = custresponse.WriteJSON(w, http.StatusAccepted, "Token was expired. New token sent to your email", nil)
			if err != nil {
				custresponse.ServerErrorResponse(w, r, err)
			}
			wg.Wait()

			return

		// Could not Deactivate Token
		case strings.Contains(err.Error(), "Could not Deactivate"):
			custresponse.ServerErrorResponse(w, r, err)
			return
		default:
			custresponse.ServerErrorResponse(w, r, err)
			return
		}
	}
	// Response Token
	err = custresponse.WriteJSON(w, http.StatusAccepted, "User activated. You may login", nil)
	if err != nil {
		custresponse.ServerErrorResponse(w, r, err)
	}

}
