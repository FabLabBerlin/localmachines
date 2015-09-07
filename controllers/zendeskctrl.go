package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/models"
)

type ZenDeskController struct {
	Controller
}

// @Title CreateTicket
// @Description Create ZenDesk ticket
// @Param name     query   string  true     "Full name of the customer"
// @Param email    query   string  true     "Email of the customer"
// @Param subject  query   string  true     "Subject of the Ticket"
// @Param message  query   string  true     "Message of the Ticket"
// @Success 200 {object} models.FastBillCreateCustomerResponse
// @Failure 500 Internal Server Error
// @Failure 401 Not authorized
// @router  /tickets [post]
func (this *ZenDeskController) CreateTicket() {
	zd := models.ZenDesk{
		Email:    beego.AppConfig.String("zendeskemail"),
		ApiToken: beego.AppConfig.String("zendeskapitoken"),
		BaseUrl:  beego.AppConfig.String("zendeskbaseurl"),
	}

	ticket := models.ZenDeskTicket{
		Requester: models.ZenDeskTicketRequester{
			Name:  this.GetString("name"),
			Email: this.GetString("email"),
		},
		Subject: this.GetString("subject"),
		Comment: models.ZenDeskTicketComment{
			Body: this.GetString("message"),
		},
	}
	if err := zd.SubmitTicket(ticket); err == nil {
		this.ServeJson()
	} else {
		beego.Error("Error submitting ZenDesk Ticket: ", err)
		this.CustomAbort(500, "Error submitting Ticket")
	}
}
