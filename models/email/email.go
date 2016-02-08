// email sending.
package email

import (
	"github.com/astaxie/beego"
	"net/smtp"
)

type Email struct {
	auth smtp.Auth
	host string
	from string
	pw   string
}

// Creates a new Email store out of the Email model.
func New() (email *Email) {
	email = &Email{}
	email.host = beego.AppConfig.String("smtphost")
	email.from = beego.AppConfig.String("emailsenderaddr")
	email.pw = beego.AppConfig.String("emailsenderpw")
	email.auth = smtp.PlainAuth("", email.from, email.pw, email.host)
	return
}

// Sends an email by consuming an existing Email store.
func (this Email) Send(to, subject, message string) error {
	msg := []byte("From: " + this.from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")
	recv := []string{to}
	err := smtp.SendMail(this.host+":25", this.auth, this.from, recv, msg)
	return err
}
