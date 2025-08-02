package usermodel

import "strings"

type TestStuctures struct {
}

func (ts *TestStuctures) NewUserOk() User {
	normalName := "TheName"
	normalMail := "mail@mail.ru"
	normalPassword := "thePassword123"

	userOk := User{
		Email: normalMail,
		Name:  normalName,
	}
	userOk.Password.Set(normalPassword)

	return userOk
}
func (ts *TestStuctures) NewUserLongName() User {
	normalMail := "mail@mail.ru"
	longName := strings.Repeat("a", 51)
	normalPassword := "thePassword123"

	userLongName := User{
		Email: normalMail,
		Name:  longName,
	}
	userLongName.Password.Set(normalPassword)

	return userLongName
}
func (ts *TestStuctures) NewUserLongMail() User {
	longName := strings.Repeat("a", 51)
	normalPassword := "thePassword123"
	tooLongMail := "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq@mail.ru"

	userLongMail := User{
		Email: tooLongMail,
		Name:  longName,
	}
	userLongMail.Password.Set(normalPassword)

	return userLongMail
}
func (ts *TestStuctures) NewUserLongPassword() User {
	normalMail := "mail@mail.ru"
	longName := strings.Repeat("a", 51)
	tooLongPassword := strings.Repeat("a", 51)

	userLongPassword := User{
		Email: normalMail,
		Name:  longName,
	}
	userLongPassword.Password.Set(tooLongPassword)

	return userLongPassword
}

func (ts *TestStuctures) NewUserInputOk() UserRegisterInput {
	return UserRegisterInput{
		Email:    "abc@mail.ru",
		Name:     "Ok name",
		Password: "ThePassword",
	}
}

func (ts *TestStuctures) NewUserInputLongName() UserRegisterInput {
	return UserRegisterInput{
		Name:     strings.Repeat("a", 51),
		Email:    "Ok name",
		Password: "ThePassword",
	}
}

func (ts *TestStuctures) NewUserInputLongMail() UserRegisterInput {
	tooLongMail := "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq@mail.ru"
	return UserRegisterInput{
		Email:    tooLongMail,
		Name:     "Ok name",
		Password: "ThePassword",
	}
}

func (ts *TestStuctures) NewUserInputTooLongPassword() UserRegisterInput {
	return UserRegisterInput{
		Email:    "abc@mail.ru",
		Name:     "Ok name",
		Password: strings.Repeat("a", 51),
	}
}
