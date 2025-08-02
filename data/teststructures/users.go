// Test structures to test users
// Could have structs with tests

package teststructures

import (
	"strings"

	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
)

func NewUserOk() usermodel.User {
	normalName := "TheName"
	normalMail := "mail@mail.ru"
	normalPassword := "thePassword123"

	userOk := usermodel.User{
		Email: normalMail,
		Name:  normalName,
	}
	userOk.Password.Set(normalPassword)

	return userOk
}
func NewUserLongName() usermodel.User {
	normalMail := "mail@mail.ru"
	longName := strings.Repeat("a", 51)
	normalPassword := "thePassword123"

	userLongName := usermodel.User{
		Email: normalMail,
		Name:  longName,
	}
	userLongName.Password.Set(normalPassword)

	return userLongName
}
func NewUserLongMail() usermodel.User {
	longName := strings.Repeat("a", 51)
	normalPassword := "thePassword123"
	tooLongMail := "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq@mail.ru"

	userLongMail := usermodel.User{
		Email: tooLongMail,
		Name:  longName,
	}
	userLongMail.Password.Set(normalPassword)

	return userLongMail
}
func NewUserLongPassword() usermodel.User {
	normalMail := "mail@mail.ru"
	longName := strings.Repeat("a", 51)
	tooLongPassword := strings.Repeat("a", 51)

	userLongPassword := usermodel.User{
		Email: normalMail,
		Name:  longName,
	}
	userLongPassword.Password.Set(tooLongPassword)

	return userLongPassword
}
