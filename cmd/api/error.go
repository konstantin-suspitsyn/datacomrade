package api

import "errors"

var ErrNoAuthorization = errors.New("Authorization header is empty")
var ErrBrokenAuthHeader = errors.New("Auth header is incorrect")
var ErrNotBearer = errors.New("Auth header is not Bearer")

var ErrAppUserNotFound = errors.New("AppUser not found")
var ErrNotAuthorization = errors.New("JWT is not authorization")
