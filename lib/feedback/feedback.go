package feedback

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/email"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
)

type Feedback struct {
	Location       *locations.Location
	RequesterName  string
	RequesterEmail string
	Subject        string
	Message        string
}

func (f *Feedback) Send() (err error) {
	if f.Location.Id <= 0 {
		return fmt.Errorf("invalid location id: %v", f.Location.Id)
	}
	if f.Location.Id == 1 {
		zd := ZenDesk{
			Email:    beego.AppConfig.String("zendeskemail"),
			ApiToken: beego.AppConfig.String("zendeskapitoken"),
			BaseUrl:  beego.AppConfig.String("zendeskbaseurl"),
		}

		ticket := ZenDeskTicket{
			Requester: ZenDeskTicketRequester{
				Name:  f.RequesterName,
				Email: f.RequesterEmail,
			},
			Subject: f.Subject,
			Comment: ZenDeskTicketComment{
				Body: f.Message,
			},
		}
		if err := zd.SubmitTicket(ticket); err != nil {
			return fmt.Errorf("Error submitting Ticket: %v", err)
		}
	} else {
		mail := email.New()
		to := f.Location.Email
		msg := "Feedback of " + f.RequesterName + ", " + f.RequesterEmail + "\n\n"
		msg += f.Message
		return mail.Send(to, f.Subject, msg)
	}
	return
}
