package bin

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/month"
)

const (
	MONTH = Unit("month")
	WEEK  = Unit("week")
	DAY   = Unit("day")
)

type Unit string

type Width struct {
	u Unit
}

func NewWidth(unit Unit) Width {
	return Width{
		u: unit,
	}
}

func (w Width) TimeFormat() string {
	switch w.u {
	case DAY:
		return "2006-01-02"
	case MONTH:
		return "2006-01"
	default:
		panic(fmt.Sprintf("unknown unit %v", w.u))
	}
}

func Map(from, to, m month.Month) (i int, ok bool) {
	for t := from; t.BeforeOrEqual(to); t = t.Add(1) {
		if t.Equal(m) {
			return i, true
		}
		i++
	}
	return -1, false
}
