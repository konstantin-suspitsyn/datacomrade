package users

import "errors"

var ErrTokenExpired = errors.New("Token is expired")
var ErrUserNotActivated = errors.New("User is not activated")
var ErrJWTSignatureFail = errors.New("Failed to sign JWT")
var ErrJWTCreateId = errors.New("Failed to create ID")
var ErrJWTIncorrectSignMethod = errors.New("Incorrect Signing Method")
var ErrJWTErrorReadingTheToken = errors.New("Could not read JWT")
var ErrJWTInvalidTokenClaims = errors.New("Invalid JWT Claims")

var ErrNotAnAccessToken = errors.New("Not and Access JWT")
var ErrNotARefreshToken = errors.New("Not and Refresh JWT")
