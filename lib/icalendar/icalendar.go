package icalendar

import (
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"time"
	"strconv"
	"strings"
)

const TIME_FMT = "20060102T150405Z"

const PRELUDE_1 = `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Makea Industries GmbH//Easylab//EN
CALSCALE:GREGORIAN
METHOD:PUBLISH
X-WR-CALNAME:reservations@fablab.berlin
X-WR-TIMEZONE:Europe/Berlin`

const PRELUDE_2 = `BEGIN:VTIMEZONE
TZID:Europe/Berlin
X-LIC-LOCATION:Europe/Berlin
BEGIN:DAYLIGHT
TZOFFSETFROM:+0100
TZOFFSETTO:+0200
TZNAME:CEST
DTSTART:19700329T020000
RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=-1SU
END:DAYLIGHT
BEGIN:STANDARD
TZOFFSETFROM:+0200
TZOFFSETTO:+0100
TZNAME:CET
DTSTART:19701025T030000
RRULE:FREQ=YEARLY;BYMONTH=10;BYDAY=-1SU
END:STANDARD
END:VTIMEZONE`

const END = `END:VCALENDAR`

type Event struct {
	Machine *machine.Machine
	Reservation *purchases.Reservation
	UserRole user_roles.Role
}

func (e Event) toIcal() (ical string) {
	ical += "BEGIN:VEVENT\n"
	ical += "DTSTART:" + e.Reservation.Purchase.TimeStart.UTC().Format(TIME_FMT) + "\n"
	ical += "DTEND:" + e.Reservation.Purchase.TimeEnd.UTC().Format(TIME_FMT) + "\n"
	ical += "DTSTAMP:" + time.Now().UTC().Format(TIME_FMT) + "\n"
	ical += "CREATED:" + e.Reservation.Purchase.Created.Format(TIME_FMT) + "\n"
	ical += "SUMMARY:" + e.Machine.Name + "\n"
	ical += "UID:mail+reservation" + strconv.FormatInt(e.Reservation.Id(), 10) + "@fablab.berlin\n"
	ical += "DESCRIPTION:Reserved by "
	if e.UserRole == user_roles.STAFF || e.UserRole == user_roles.ADMIN || e.UserRole == user_roles.SUPER_ADMIN {
		ical += "Staff\n"
	} else {
		ical += "Member\n"
	}
	ical += "SEQUENCE:0\n"
	ical += "STATUS:"
	if e.Reservation.Purchase.Cancelled {
		ical += "CANCELLED\n"
	} else if e.Reservation.Purchase.ReservationDisabled {
		ical += "CANCELLED\n"
	} else {
		ical += "CONFIRMED\n"
	}
	ical += "END:VEVENT\n"
	return
}

func ToIcal(es []Event) (ical string) {
	ical += PRELUDE_1 + "\n"
	//ical += PRELUDE_2 + "\n"
	for _, e := range es {
		ical += e.toIcal()
	}
	ical += END + "\n"
	ical = strings.Replace(ical, "\n", "\r\n", -1)
	return
}
