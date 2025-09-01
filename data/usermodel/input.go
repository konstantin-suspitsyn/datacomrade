package usermodel

// Input json structure for user to register
type UserRegisterInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailInput struct {
	Email string `json:"email"`
}

type PasswordInput struct {
	Password string `json:"password"`
}

type UpdatePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type RefreshTokenUpdater struct {
	RefreshToken string `json:"refresh_token"`
}
