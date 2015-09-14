package models

import (
	"github.com/astaxie/beego"
	"net/smtp"
)

type Email struct {
	auth smtp.Auth
	host string
	from string
}

func NewEmail() (this Email) {
	this.host = beego.AppConfig.String("smtphost")
	this.from = beego.AppConfig.String("emailsenderaddr")
	this.pw = beego.AppConfig.String("emailsenderpw")
	this.auth = smtp.PlainAuth("", this.from, pw, this.host)
	return
}

func (this Email) Send(to, subject, message string) error {
	beego.Info("Email#Send: from:", this.from)
	beego.Info("Email#Send: to:", to)
	beego.Info("Email#Send: subject:", subject)
	beego.Info("Email#Send: message:", message)
	msg := []byte("From: " + this.from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")
	recv := []string{to}
	err := smtp.SendMail(this.host+":25", this.auth, this.from, recv, msg)
	return err
}
