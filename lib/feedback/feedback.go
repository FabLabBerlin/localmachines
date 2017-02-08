package feedback

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/email"
	"github.com/FabLabBerlin/localmachines/models/locations"
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

	mail := email.New()
	to := f.Location.Email
	msg := "Feedback of " + f.RequesterName + ", " + f.RequesterEmail + "\n\n"
	msg += f.Message
	return mail.Send(to, f.Subject, msg)

	return
}
