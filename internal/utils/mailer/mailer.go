package mailer

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"gopkg.in/gomail.v2"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	sender string
}

func NewMailer() Mailer {

	mailConfig := configs.InitMailConfig()

	mail := New(mailConfig.MAIL_HOST, mailConfig.MAIL_PORT, mailConfig.MAIL_USER, mailConfig.MAIL_PASSWORD, mailConfig.MAIL_EMAIL)

	return mail

}

func New(host string, port int, username string, password string, sender string) Mailer {
	dialer := gomail.NewDialer(host, port, username, password)

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("subject", subject.String())

	msg.SetBody("text/plain", plainBody.String())
	msg.SetBody("text/html", htmlBody.String())

	// TODO: MAYBE TRY SEND MAIL MULTIPLE TIMES
	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
