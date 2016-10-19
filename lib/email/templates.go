package email

import (
	"text/template"
)

const locationReservationNotification = `
Hello {{.Location.Title}} Team!

The user {{.User.FirstName}} {{.User.LastName}} ({{.User.Email}}) has made
a reservation for {{.Machine.Name}} on

{{.TimeStart}} lasting for {{.Reservation.Purchase.Duration}}

Cheers

EASY LAB
`

const userReservationNotification = `
Hello {{.User.FirstName}}!

You made a reservation for {{.Machine.Name}} on

{{.TimeStart}} lasting for {{.Reservation.Purchase.Duration}}

Happy Making!
`

var (
	LocationReservationNotification *template.Template
	UserReservationNotification     *template.Template
)

func init() {
	var err error

	LocationReservationNotification = template.New("LocationReservationNotification")
	LocationReservationNotification, err = LocationReservationNotification.Parse(locationReservationNotification)
	if err != nil {
		panic(err.Error())
	}

	UserReservationNotification = template.New("UserReservationNotification")
	UserReservationNotification, err = UserReservationNotification.Parse(userReservationNotification)
	if err != nil {
		panic(err.Error())
	}
}
