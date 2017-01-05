package bin

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"time"
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

func (w Width) IsMonth() bool {
	return w.u == MONTH
}

func (w Width) TimeIndex(t time.Time) string {
	switch w.u {
	case DAY:
		return t.Format("2006-01-02")
	case WEEK:
		y, w := t.ISOWeek()
		return fmt.Sprintf("%v/%0.2d", y, w)
	case MONTH:
		return t.Format("2006-01")
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
