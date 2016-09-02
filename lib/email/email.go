// email sending.
package email

import (
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego"
	"net/smtp"
)

type Email struct {
	auth      smtp.Auth
	host      string
	from      string
	pw        string
	tlsConfig *tls.Config
}

// Creates a new Email store out of the Email model.
func New() (email *Email) {
	email = &Email{
		host: beego.AppConfig.String("smtphost"),
		from: beego.AppConfig.String("emailsenderaddr"),
		pw:   beego.AppConfig.String("emailsenderpw"),
	}
	email.auth = smtp.PlainAuth("", email.from, email.pw, email.host)
	email.tlsConfig = &tls.Config{
		ServerName: email.host,
	}
	return
}

// Sends an email by consuming an existing Email store.
func (this Email) Send(to, subject, message string) error {
	msg := []byte("From: " + this.from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")
	//recv := []string{to}

	//err := smtp.SendMail(this.host+":25", this.auth, this.from, recv, msg)
	err := this.sendTlsMail(this.from, to, msg)
	return err
}

// https://gist.github.com/chrisgillis/10888032
func (this Email) sendTlsMail(from, to string, msg []byte) (err error) {
	conn, err := tls.Dial("tcp", this.host+":465", this.tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial: %v", err)
	}
	c, err := smtp.NewClient(conn, this.host)
	if err != nil {
		return
	}

	if err = c.Auth(this.auth); err != nil {
		return
	}

	if err = c.Mail(from); err != nil {
		return
	}

	if err = c.Rcpt(to); err != nil {
		return
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return
	}

	_, err = w.Write(msg)
	if err != nil {
		return
	}

	err = w.Close()
	if err != nil {
		return
	}

	c.Quit()
	return
}
