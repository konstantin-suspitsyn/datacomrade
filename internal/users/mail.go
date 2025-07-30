package users

import (
	"log/slog"

	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/mailer"
)

func (us *UserService) SendActivationToken(tokenPlainText, email string) {

	mailer := mailer.NewMailer()

	data := map[string]any{"activationToken": tokenPlainText}

	err := mailer.Send(email, "user_welcome.tmpl", data)

	if err != nil {
		slog.Info(err.Error(), "email", email)
	}
}

func (us *UserService) SendForgotPasswordMail(tokenPlainText, email string) {
	mailer := mailer.NewMailer()

	data := map[string]any{"forgotPassword": tokenPlainText}
	err := mailer.Send(email, "forgot_password.tmpl", data)

	if err != nil {
		slog.Info(err.Error(), "email", email)
	}
}
