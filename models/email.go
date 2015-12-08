package models

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
func NewEmail() (this Email) {
	this.host = beego.AppConfig.String("smtphost")
	this.from = beego.AppConfig.String("emailsenderaddr")
	this.pw = beego.AppConfig.String("emailsenderpw")
	this.auth = smtp.PlainAuth("", this.from, this.pw, this.host)
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
