package day

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"strconv"
	"strings"
	"time"
)

// Day representation independent of time (zone).
type Day struct {
	m month.Month
	d int
}

func New(y int, m time.Month, d int) Day {
	if d > 31 {
		panic(fmt.Sprintf("%v"))
	}
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

func Now() (d Day) {
	return NewTime(time.Now())
}

func (d Day) Add(dur time.Duration) Day {
	t := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	return NewTime(t.Add(dur))
}

func (d Day) AddDate(y, m, dd int) Day {
	t := time.Date(d.Year(), d.Month(), d.Day(), 11, 0, 0, 0, time.UTC)
	t = t.AddDate(y, m, dd)

	return NewTime(t)
}

// addMonth where the neutral element of addition is the last day of month.
func (d Day) addMonth() Day {
	s := time.Date(d.Year(), d.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	u := time.Date(d.Year(), d.Month()+2, 1, 0, 0, 0, 0, time.UTC)
	return d.Add(u.Sub(s))
}

// addMonths where the neutral element of addition is the last day of month.
func (d Day) addMonths(months int) (e Day) {
	e = d

	for i := 0; i < months; i++ {
		e = e.addMonth()
	}

	return
}

// AddDate2 where the neutral element of monthly addition is the last day of
// month.
//
// Example: day.New(2016, 2, 29).AddDate2(0, 1, 0) => March 31st
//
//           but
//
//          day.New(2016, 2, 29).AddDate(0, 1, 0) => March 29th
func (d Day) AddDate2(y, m, dd int) Day {
	return d.addMonths(m).AddDate(y, 0, dd)
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

func (d Day) BeforeOrEqual(other Day) bool {
	return d.Before(other) || d.Equal(other)
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

func (d Day) IsZero() bool {
	return d.d == 0 && d.m.IsZero()
}

func (d Day) Month() time.Month {
	return d.m.Month()
}

func (d Day) String() string {
	t := time.Date(d.Year(), d.Month(), d.Day(), 11, 0, 0, 0, time.UTC)
	return t.Format("2006-01-02")
}

func (d Day) Sub(other Day) time.Duration {
	t := time.Date(d.Year(), d.Month(), d.Day(), 11, 0, 0, 0, time.UTC)
	u := time.Date(other.Year(), other.Month(), other.Day(), 11, 0, 0, 0, time.UTC)
	return t.Sub(u)
}

func (d Day) Year() int {
	return d.m.Year()
}
