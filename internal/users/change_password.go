package users

import "github.com/konstantin-suspitsyn/datacomrade/data/usermodel"

func (us *UserService) ChangeOldPassword(userId int64, oldPassword, newPassword string) error {

	user, err := us.Models.User.GetById(userId)
	if err != nil {
		return err
	}
	passwordMatch, err := user.Password.Matches(oldPassword)

	if err != nil {
		return err
	}

	if passwordMatch {

		us.Models.User.UpdatePassword(userId, newPassword)
		return nil
	}

	return usermodel.ErrPasswordsDoNotMatch

}
