package users

// Checks email and password against database and returns true if user exists
// returns false and error if something wrong with user
func (us *UserService) checkEmailPassword(email string, password string) (bool, error) {

	user, err := us.Models.User.GetByEmail(email)
	if err != nil {
		return false, err
	}

	// Check if user is activated
	if !user.Activated {
		return false, ErrUserNotActivated
	}

	passwordMatch, err := user.Password.Matches(password)

	if err != nil {
		return false, err
	}

	return passwordMatch, nil
}
