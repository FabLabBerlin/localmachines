package day

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"strconv"
	"strings"
	"time"
)

type Day struct {
	m month.Month
	d int
}

func New(d int, m time.Month, y int) Day {
	return Day{
		d: d,
		m: month.New(m, y),
	}
}

// NewString("2015-01-15") results in Jan 15th, 2015
func NewString(s string) (d Day, err error) {
	if tmp := strings.Split(s, "-"); len(tmp) == 3 {
		d.d, err = strconv.Atoi(tmp[2])
		if err != nil {
			return d, fmt.Errorf("day format: %v", err)
		}

		s := strings.Join(tmp[:2], "-")
		d.m, err = month.NewString(s)
		if err != nil {
			return d, fmt.Errorf("month/year format: %v", err)
		}

		return d, err
	} else {
		return d, errors.New("wrong format")
	}
}

func NewTime(t time.Time) (d Day) {
	return Day{
		d: t.Day(),
		m: month.New(t.Month(), t.Year()),
	}
}

func (d Day) After(other Day) bool {
	if d.m.Equal(other.m) {
		return d.d > other.d
	} else {
		return d.m.After(other.m)
	}
}

func (d Day) AfterTime(t time.Time) bool {
	return d.After(NewTime(t))
}

func (d Day) Before(other Day) bool {
	if d.m.Equal(other.m) {
		return d.d < other.d
	} else {
		return d.m.Before(other.m)
	}
}

func (d Day) BeforeTime(t time.Time) bool {
	return d.Before(NewTime(t))
}

func (d Day) Contains(t time.Time, loc *time.Location) bool {
	if !d.m.Contains(t, loc) {
		return false
	}

	return t.In(loc).Day() == d.d
}

func (d Day) Day() int {
	return d.d
}

func (d Day) Equal(other Day) bool {
	return d.m.Equal(other.m) && d.d == other.d
}

func (d Day) Month() time.Month {
	return d.m.Month()
}

func (d Day) String() string {
	t := time.Date(d.Year(), d.Month(), d.Day(), 11, 0, 0, 0, time.UTC)
	return t.Format("2006-01-02")
}

func (d Day) Year() int {
	return d.m.Year()
}
