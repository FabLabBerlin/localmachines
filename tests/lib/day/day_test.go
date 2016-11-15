package day

import (
	"github.com/FabLabBerlin/localmachines/lib/day"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestLibDay(t *testing.T) {
	Convey("Construct and Copy", t, func() {
		d := day.New(7, time.February, 2014)
		var cpy day.Day
		cpy = d
		So(cpy.Day(), ShouldEqual, 7)
		So(cpy.Month(), ShouldEqual, time.February)
		So(cpy.Year(), ShouldEqual, 2014)
	})

	Convey("NewString", t, func() {
		m, err := day.NewString("2013-07-16")
		So(err, ShouldBeNil)
		So(m.Day(), ShouldEqual, 16)
		So(m.Month(), ShouldEqual, time.July)
		So(m.Year(), ShouldEqual, 2013)
	})

	Convey("After", t, func() {
		d := day.New(28, time.February, 2014)
		before := day.New(27, time.February, 2014)
		after := day.New(11, time.January, 2015)

		So(before.Before(d), ShouldBeTrue)
		So(after.After(d), ShouldBeTrue)
		So(d.Before(d), ShouldBeFalse)
	})

	Convey("Contains", t, func() {
		loc, err := time.LoadLocation("Europe/Berlin")
		if err != nil {
			panic(err.Error())
		}

		tm := time.Date(2015, 11, 1, 0, 0, 0, 0, loc)
		d := day.New(1, time.November, 2015)
		So(d.Contains(tm, loc), ShouldBeTrue)
	})

	Convey("Equal", t, func() {
		d := day.New(7, time.November, 2015)
		e := day.New(8, time.November, 2015)
		So(d.Equal(e), ShouldBeFalse)
		So(d.Equal(d), ShouldBeTrue)
	})

	Convey("String", t, func() {
		m := day.New(17, time.November, 2015)
		So(m.String(), ShouldEqual, "2015-11-17")
	})
}
