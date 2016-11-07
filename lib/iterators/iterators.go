package iterators

import (
	"time"
)

type Monthly struct {
	start Month
	end   Month

	current *Month
}

type Month struct {
	Year  int
	Month time.Time
}

func (m Month) Next() Month {
	y := m.Year
	m := int(m.Month) + 1

	if m > 12 {
		m = 1
		y++
	}

	return Month{
		Year:  y,
		Month: m,
	}
}

func NewMonthly(start Month, end Month) *Monthly {
	return &Monthly{
		start: start,
		end:   end,
	}
}

func (m *Monthly) Begin() *Month {
	current := m.start
	m.current = &current

	return m.current
}

func (m *Monthly) Next() *Month {
	if m.current != nil {
		n := m.current.Next()
		if n.LessThan(m.end) {
			m.current = &n
		} else {
			m.current = nil
		}
	}

	return nil
}
