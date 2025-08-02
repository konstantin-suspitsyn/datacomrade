package users

import (
	"testing"

	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/comradetest"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	db, err := comradetest.CreateDbMock()
	if err != nil {
		t.Errorf("Did not create DB")
	}
	defer db.Close()

	comradetest.InitEnv()

	userService := New(db)

	registerInputOk := usermodel.UserRegisterInput{
		Email:    "abc@mail.ru",
		Name:     "Ok name",
		Password: "ThePassword",
	}
	user, err, _ := userService.createUser(registerInputOk)

	assert.Equal(t, registerInputOk.Email, user.Email, "Checking user email")
	assert.Equal(t, registerInputOk.Email, user.Email, "Checking user email")
	matches, err := user.Password.Matches(registerInputOk.Password)
	if err != nil {
		t.Errorf("Match error failed")
	}
	assert.Equal(t, true, matches, "Checking user email")
}
