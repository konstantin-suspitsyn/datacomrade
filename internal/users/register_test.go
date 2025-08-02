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

	testData := usermodel.TestStuctures{}

	testIt := []struct {
		UserInput usermodel.UserRegisterInput
		HasError  bool
		Condition string
	}{
		{UserInput: testData.NewUserInputOk(), HasError: false, Condition: "Fine User"},
		{UserInput: testData.NewUserInputLongName(), HasError: true, Condition: "Too long name"},
		{UserInput: testData.NewUserInputLongMail(), HasError: true, Condition: "Too long mail"},
		{UserInput: testData.NewUserInputTooLongPassword(), HasError: true, Condition: "Too long Password"},
	}

	for _, tt := range testIt {
		user, err, _ := userService.createUser(tt.UserInput)
		if err != nil {
			assert.Equal(t, tt.HasError, true, tt.Condition)
		} else {
			assert.Equal(t, tt.UserInput.Email, user.Email, "Checking user email")
			assert.Equal(t, tt.UserInput.Email, user.Email, "Checking user email")
			matches, errPass := user.Password.Matches(tt.UserInput.Password)
			if errPass != nil {
				t.Errorf("Match error failed")
			}
			assert.Equal(t, true, matches, "Checking user email", tt.Condition)

		}

	}
}
