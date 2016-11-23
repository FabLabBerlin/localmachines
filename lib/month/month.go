package month

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

// NewString("2015-01") or NewString("2015-01-15") results in Jan '15
func NewString(s string) (m Month, err error) {
	tmp := strings.Split(s, "-")

	switch len(tmp) {
	case 2, 3:
		m.y, err = strconv.Atoi(tmp[0])
		if err != nil {
			return m, fmt.Errorf("year format: %v", err)
		}
		mm, err := strconv.Atoi(tmp[1])
		if err != nil {
			return m, fmt.Errorf("month format: %v", err)
		}
		m.m = time.Month(mm)

		return m, err
	default:
		return m, errors.New("wrong format")
	}
}

func NewTime(t time.Time) (m Month) {
	return Month{
		y: t.Year(),
		m: t.Month(),
	}
}

func (m Month) After(other Month) bool {
	if m.Year() == other.Year() {
		return int(m.Month()) > int(other.Month())
	} else {
		return m.Year() > other.Year()
	}
}

func (m Month) AfterOrEqual(other Month) bool {
	return m.After(other) || m.Equal(other)
}

func (m Month) AfterTime(t time.Time) bool {
	return m.After(NewTime(t))
}

func (m Month) Before(other Month) bool {
	if m.Year() == other.Year() {
		return int(m.Month()) < int(other.Month())
	} else {
		return m.Year() < other.Year()
	}
}

func (m Month) BeforeOrEqual(other Month) bool {
	return m.Before(other) || m.Equal(other)
}

func (m Month) BeforeTime(t time.Time) bool {
	return m.Before(NewTime(t))
}

func (m Month) Contains(t time.Time, loc *time.Location) bool {
	u := t.In(loc)
	return u.Month() == m.Month() && u.Year() == m.Year()
}

func (m Month) Equal(other Month) bool {
	return m.Month() == other.Month() && m.Year() == other.Year()
}

func (m Month) Month() time.Month {
	return time.Month(m.m)
}

func (m Month) String() string {
	t := time.Date(m.Year(), m.Month(), 11, 0, 0, 0, 0, time.UTC)
	return t.Format("2006-01")
}

func (m Month) Year() int {
	return m.y
}
