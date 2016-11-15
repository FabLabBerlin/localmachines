package month

import (
	"time"
)

type Month struct {
	m time.Month 
	y int 
}

func New(m time.Month, y int) Month {
	return Month{
		m: m,
		y: y,
	}
}

func (m Month) After(other Month) bool {
	if m.Year() == other.Year() {
		return int(m.Month()) > int(other.Month())
	} else {
		return m.Year() > other.Year()
	}
}

func (m Month) Before(other Month) bool {
	if m.Year() == other.Year() {
		return int(m.Month()) < int(other.Month())
	} else {
		return m.Year() < other.Year()
	}
}

func (m Month) Contains(t time.Time, loc *time.Location) bool {
	u := t.In(loc)
	return u.Month() == m.Month() && u.Year() == m.Year()
}

func (m Month) Equal(other Month) bool{
	return m.Month() == other.Month() && m.Year() == other.Year()
}

func (m Month) Month() time.Month {
	return time.Month(m.m)
}

func (m Month) Year() int {
	return m.y
}
