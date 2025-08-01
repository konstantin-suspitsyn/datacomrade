package usermodel

import (
	"log/slog"
	"testing"

	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/validator"
	"github.com/stretchr/testify/assert"
)

func (suite *UserModelSuite) TestOne() {
	t := suite.T()

	assert.Equal(t, 1, 1, "The result should be equal to the expected value")
}

func TestValidateUser(t *testing.T) {
	teststructures := TestStuctures{}
	userOk := teststructures.NewUserOk()
	userLongName := teststructures.NewUserLongName()
	userLongMail := teststructures.NewUserLongMail()
	userLongPassword := teststructures.NewUserLongPassword()

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
		err := suite.Model.User.Insert(suite.ctx, &tt.User)
		assert.Equal(t, tt.err, err, tt.Condition)
	}
}

func (suite *UserModelSuite) TestGetByEmail() {
	t := suite.T()
	name := "TestGetByEmailUser"
	mail := "testbymail@mail.ru"
	password := "thePassword123"

	userNormal := User{
		Name:  name,
		Email: mail,
	}
	userNormal.Password.Set(password)

	err := suite.Model.User.Insert(suite.ctx, &userNormal)
	if err != nil {
		t.Errorf("Insert is broken")
	}

	user, err := suite.Model.User.GetByEmail(suite.ctx, mail)

	assert.Equal(t, userNormal.Name, user.Name, "Get user by email")

}

func (suite *UserModelSuite) TestGetById() {
	t := suite.T()
	name := "TestById"
	mail := "testbyid@mail.ru"
	password := "thePassword123"

	userTestById := User{
		Name:  name,
		Email: mail,
	}
	userTestById.Password.Set(password)

	err := suite.Model.User.Insert(suite.ctx, &userTestById)
	if err != nil {
		t.Errorf("Insert is broken")
	}

	user, err := suite.Model.User.GetById(suite.ctx, userTestById.Id)

	assert.Equal(t, userTestById.Name, user.Name, "Get user by id")

}

func (suite *UserModelSuite) TestUpdatePassword() {
	t := suite.T()
	name := "UpdatePassUser"
	mail := "testbypass@mail.ru"
	password := "thePassword123"
	passwordNew := "thePassword1235"

	userPaswordChange := User{
		Name:  name,
		Email: mail,
	}
	userPaswordChange.Password.Set(password)

	err := suite.Model.User.Insert(suite.ctx, &userPaswordChange)
	if err != nil {
		slog.Info(err.Error())
		t.Errorf("Insert is broken")
	}

	err = suite.Model.User.UpdatePassword(suite.ctx, userPaswordChange.Id, passwordNew)

	if err != nil {
		slog.Info(err.Error())
		t.Errorf("Password change update broken")
	}

	user, err := suite.Model.User.GetById(suite.ctx, userPaswordChange.Id)

	if err != nil {
		slog.Info(err.Error())
		t.Errorf("Did not find user")
	}

	match, err := user.Password.Matches(passwordNew)

	if err != nil {
		slog.Info(err.Error())
		t.Errorf("Error matching password")
	}

	assert.Equal(t, true, match, "Match passwords")

}
