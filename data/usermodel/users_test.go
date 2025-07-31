package usermodel

import (
	"strings"
	"testing"

	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/validator"
	"github.com/stretchr/testify/assert"
)

func (suite *UserModelSuite) TestOne() {
	t := suite.T()

	assert.Equal(t, 1, 1, "The result should be equal to the expected value")
}

func TestValidateUser(t *testing.T) {

	normalName := "TheName"
	normalMail := "mail@mail.ru"
	longName := strings.Repeat("a", 51)
	normalPassword := "thePassword123"
	tooLongPassword := strings.Repeat("a", 51)
	tooLongMail := "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq@mail.ru"

	userOk := User{
		Email: normalMail,
		Name:  normalName,
	}
	userOk.Password.Set(normalPassword)

	userLongName := User{
		Email: normalMail,
		Name:  longName,
	}
	userLongName.Password.Set(normalPassword)

	userLongMail := User{
		Email: tooLongMail,
		Name:  longName,
	}
	userLongMail.Password.Set(normalPassword)

	userLongPassword := User{
		Email: normalMail,
		Name:  longName,
	}
	userLongPassword.Password.Set(tooLongPassword)

	tests := []struct {
		condition string
		user      User
		hasError  bool
	}{
		{condition: "Ok", user: userOk, hasError: false},
		{condition: "LongName", user: userLongName, hasError: true},
		{condition: "LongMail", user: userLongMail, hasError: true},
		{condition: "LongPassword", user: userLongPassword, hasError: true},
	}

	for _, tt := range tests {
		v := validator.New()
		ValidateUser(v, &tt.user)
		assert.Equal(t, !tt.hasError, v.Valid(), tt.condition)
	}

}

func (suite *UserModelSuite) TestInsert() {

	t := suite.T()
	normalName := "TheName"
	normalNameMore := "TheName1"
	normalMail := "mail@mail.ru"
	normalMailMore := "mail1@mail.ru"
	normalPassword := "thePassword123"

	userNormal := User{
		Name:  normalName,
		Email: normalMail,
	}
	userNormal.Password.Set(normalPassword)

	userSameName := User{
		Name:  normalName,
		Email: normalMailMore,
	}
	userSameName.Password.Set(normalPassword)
	userSameEmail := User{
		Name:  normalNameMore,
		Email: normalMail,
	}
	userSameEmail.Password.Set(normalPassword)

	tests := []struct {
		Condition string
		User      User
		err       error
	}{
		{Condition: "Create first user", User: userNormal, err: nil},
		{Condition: "Create user with duplicate Name", User: userSameName, err: ErrDuplicateName},
		{Condition: "Create user with duplicate Email", User: userSameEmail, err: ErrDuplicateEmail},
	}

	for _, tt := range tests {
		err := suite.Model.User.Insert(&tt.User)
		assert.Equal(t, tt.err, err, tt.Condition)
	}
}
