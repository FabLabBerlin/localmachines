package bin

import (
	"fmt"
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
