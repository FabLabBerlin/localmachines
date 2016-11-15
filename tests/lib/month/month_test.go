package month

import (
	"github.com/FabLabBerlin/localmachines/lib/month"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestLibMonth(t *testing.T) {
	Convey("Construct and Copy", t, func() {
		m := month.New(time.February, 2014)
		var cpy month.Month
		cpy = m
		So(cpy.Month(), ShouldEqual, time.February)
		So(cpy.Year(), ShouldEqual, 2014)
	})

	Convey("After", t, func() {
		m := month.New(time.February, 2014)
		before := month.New(time.January, 2014)
		after := month.New(time.January, 2015)

		So(before.Before(m), ShouldBeTrue)
		So(after.After(m), ShouldBeTrue)
		So(m.Before(m), ShouldBeFalse)
	})

	Convey("Contains", t, func() {
		loc, err := time.LoadLocation("Europe/Berlin")
		if err != nil {
			panic(err.Error())
		}

		tm := time.Date(2015, 11, 1, 0, 0, 0, 0, loc)
		m := month.New(time.November, 2015)
		So(m.Contains(tm, loc), ShouldBeTrue)
	})

	Convey("Equal", t, func() {
		m := month.New(time.November, 2015)
		n := month.New(time.December, 2015)
		So(m.Equal(n), ShouldBeFalse)
		So(m.Equal(m), ShouldBeTrue)
	})
}
