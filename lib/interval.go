package lib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const ISO8601_FMT = "%d-%02d-%02d"

type Interval struct {
	MonthFrom int
	YearFrom  int
	MonthTo   int
	YearTo    int
}

func NewInterval(from, to string) (i Interval, err error) {
	if i.MonthFrom, i.YearFrom, err = parseYYYYMM(from); err != nil {
		return i, fmt.Errorf("from:", err)
	}
	if i.MonthTo, i.YearTo, err = parseYYYYMM(to); err != nil {
		return i, fmt.Errorf("to:", err)
	}
	return
}

func parseYYYYMM(s string) (month, year int, err error) {
	tmp := strings.Split(s, "-")
	if len(tmp) != 2 {
		return 0, 0, fmt.Errorf("wrong fmt: %v", s)
	}
	if month, err = strconv.Atoi(tmp[1]); err != nil {
		return 0, 0, fmt.Errorf("month: %v", err)
	}
	if year, err = strconv.Atoi(tmp[0]); err != nil {
		return 0, 0, fmt.Errorf("year: %v", err)
	}
	return
}

func (i Interval) Contains(t time.Time) bool {
	return i.TimeFrom().Before(t) && i.TimeTo().After(t)
}

func (i Interval) DayFrom() string {
	return fmt.Sprintf(ISO8601_FMT, i.YearFrom, i.MonthFrom, 1)
}

func (i Interval) DayTo() string {
	t := time.Date(i.YearTo, time.Month(i.MonthTo), 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 1, -1)
	return t.Format("2006-01-02")
}

func (i Interval) TimeFrom() time.Time {
	return time.Date(i.YearFrom, time.Month(i.MonthFrom), 1, 0, 0, 0, 0, time.UTC)
}

func (i Interval) TimeTo() time.Time {
	t := time.Date(i.YearTo, time.Month(i.MonthTo), 1, 23, 59, 59, 0, time.UTC)
	t = t.AddDate(0, 1, -1)
	return t
}

func (i Interval) OneMonth() bool {
	return i.MonthFrom == i.MonthTo && i.YearFrom == i.YearTo
}

func (i Interval) String() (s string) {
	s = strconv.Itoa(i.MonthFrom) + "-" + strconv.Itoa(i.YearFrom)
	if !i.OneMonth() {
		s += "-" + strconv.Itoa(i.MonthTo) + "-" + strconv.Itoa(i.YearTo)
	}
	return
}
