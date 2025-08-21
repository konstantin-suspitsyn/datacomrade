package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/konstantin-suspitsyn/datacomrade/data/shareddata"
	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
	"github.com/konstantin-suspitsyn/datacomrade/internal/services"
	"github.com/konstantin-suspitsyn/datacomrade/internal/users"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/custresponse"
)

func GetAuthMiddlewareFunc(services *services.ServiceLayer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var appUser usermodel.AppUser
			// Read auth header
			// Verify token
			userClaims, err := GetUserFromBearer(r, services)
			if err != nil {
				switch {
				case errors.Is(err, ErrNoAuthorization):
					slog.Info("EMPTY USER")
					appUser = usermodel.AppUser{}
				default:
					custresponse.ServerErrorResponse(w, r, err)
					return
				}
			} else {
				// Pass payload to context
				appUser = usermodel.AppUser{
					Id:       userClaims.Id,
					UserName: userClaims.Subject,
				}

				slog.Info("AppUser", "User name", appUser.UserName)
			}

			ctx := context.WithValue(r.Context(), shareddata.AuthKey{}, &appUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromBearer(r *http.Request, services *services.ServiceLayer) (*users.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, ErrNoAuthorization
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return nil, ErrBrokenAuthHeader
	}
	if fields[0] != "Bearer" {
		return nil, ErrNotBearer
	}

	token := fields[1]

	userClaims, err := services.UserService.JWTMaker.VerifyAccessToken(token)

	if err != nil {
		return nil, err
	}
	if userClaims.Scope != users.ScopeAuthToken {
		return nil, ErrNotAuthorization
	}

	return userClaims, nil
}

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appUser usermodel.AppUser

		if user, ok := r.Context().Value(shareddata.AuthKey{}).(*usermodel.AppUser); !ok {
			custresponse.UnauthorizedResponse(w, r)
			return
		} else {
			appUser = *user
		}

		if appUser.Id == 0 {
			custresponse.UnauthorizedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)

	})
}
