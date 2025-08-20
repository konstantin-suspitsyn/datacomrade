package users

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/konstantin-suspitsyn/datacomrade/data/shareddata"
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
	token, err := us.createRegistrationToken(ctx, &user)

	// Send email in background
	var wg sync.WaitGroup

	shared.BackroundJob(func() {
		us.SendActivationToken(token.PlainText, user.Email)
	}, &wg)

	if err != nil {
		custresponse.ServerErrorResponse(w, r, err)

		return
	}

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

	ctx := r.Context()

	token, err := us.UserModels.Token.GetByPlainText(ctx, input.TokenPlainText, usermodel.ScopeActivation)

	if err != nil {
		custresponse.ServerErrorResponse(w, r, fmt.Errorf("ERROR: %w. Getting Token from input", err))
		return
	}
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

func (us *UserService) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get login and password
	var userLoginInput usermodel.UserLoginInput

	err, validationProblems := us.ReadLogin(w, r, &userLoginInput)
	if err != nil {
		switch {
		case errors.Is(err, usermodel.ErrLoginInputValidation):
			custresponse.FailedValidationResponse(w, r, err, validationProblems)
			return
		default:
			custresponse.ServerErrorResponse(w, r, err)
			return
		}
	}

	isOk, err := us.checkEmailPasswordAndActive(ctx, userLoginInput.Email, userLoginInput.Password)

	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotActivated):
			custresponse.WriteJSON(w, http.StatusForbidden, "User is not activated", nil)
			return
		default:
			custresponse.ServerErrorResponse(w, r, err)
			return
		}
	}

	if !isOk {
		custresponse.InvalidCredentialsResponse(w, r)
		return
	}

	user, roles, err := us.findUserAndJWTRoles(ctx, userLoginInput.Email)

	if err != nil {
		custresponse.ServerErrorResponse(w, r, err)
		return
	}

	loginDTO, err := us.accessAndRefreshTokens(ctx, user, roles)

	err = custresponse.WriteJSON(w, http.StatusOK, loginDTO, nil)
}

func (us *UserService) UserForgotPassword(w http.ResponseWriter, r *http.Request) {
	var emailInput usermodel.EmailInput
	ctx := r.Context()

	err := custresponse.ReadJSON(w, r, &emailInput)

	if err != nil {
		custresponse.BadRequestResponse(w, r, fmt.Errorf("Error reading email. %w", err))
		return
	}

	user, err := us.GetUserByEmail(ctx, emailInput.Email)

	if err != nil {
		custresponse.BadRequestResponse(w, r, fmt.Errorf("Error finding user. %w", err))
		return
	}

	token, err := us.CreateForgotPasswordToken(ctx, user)
	// Send email in background
	var wg sync.WaitGroup

	shared.BackroundJob(func() {
		us.SendForgotPasswordMail(user.Email, token.PlainText)
	}, &wg)

	answer := map[string]string{"message": "Forgot password token sent"}

	custresponse.WriteJSON(w, http.StatusOK, answer, nil)

	wg.Wait()

}

func (us *UserService) ChangeForgotternPassword(w http.ResponseWriter, r *http.Request) {
	refreshToken := chi.URLParam(r, "refresh")
	var passwordInput usermodel.PasswordInput

	ctx := r.Context()
	// Look for refresh token
	token, err := us.UserModels.Token.GetByPlainText(ctx, refreshToken, usermodel.ScopeForgotPassword)

	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}
	// Validate password
	err = custresponse.ReadJSON(w, r, &refreshToken)
	// TODO: VALIDATE PASSWORD

	// Change password
	err = us.UserModels.User.UpdatePassword(ctx, token.UserId, passwordInput.Password)
	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}
}

func (us *UserService) Me(w http.ResponseWriter, r *http.Request) {

	appUser := r.Context().Value(shareddata.AuthKey{}).(*usermodel.AppUser)

	custresponse.WriteJSON(w, http.StatusOK, appUser, nil)
}
