package usermodel

import (
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/data/dataerrors"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/validator"
)

var ErrTokenRecordNotFound = fmt.Errorf("Token. %w", dataerrors.ErrRecordNotFound)
var ErrUserRecordNotFound = fmt.Errorf("Users. %w", dataerrors.ErrRecordNotFound)

var ErrRolesRecordNotFound = fmt.Errorf("Roles. %w", dataerrors.ErrRecordNotFound)

var ErrDuplicateEmail = errors.New("Duplicate email")
var ErrDuplicateName = errors.New("Duplicate name")

var ErrTokenInputValidation = fmt.Errorf("ERROR: TokenInput validation. %w", validator.ErrValidation)
var ErrLoginInputValidation = fmt.Errorf("ERROR: Email validation. %w", validator.ErrValidation)

var ErrNoRefreshTokenRecord = fmt.Errorf("ERROR: Refresh token was not found")

var ErrEmailValidation = fmt.Errorf("ERROR: email validation failed")

var ErrPasswordsDoNotMatch = fmt.Errorf("ERROR: passwords do not match")
