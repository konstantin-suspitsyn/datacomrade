package users

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
)

func (us *UserService) ChangeOldPassword(ctx context.Context, userId int64, oldPassword, newPassword string) error {

	user, err := us.Models.User.GetById(ctx, userId)
	if err != nil {
		return err
	}
	passwordMatch, err := user.Password.Matches(oldPassword)

	if err != nil {
		return err
	}

	if passwordMatch {

		us.Models.User.UpdatePassword(ctx, userId, newPassword)
		return nil
	}

	return usermodel.ErrPasswordsDoNotMatch

}
