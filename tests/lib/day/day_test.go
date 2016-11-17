package day

import (
	"github.com/FabLabBerlin/localmachines/lib/day"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestLibDay(t *testing.T) {
	Convey("Construct and Copy", t, func() {
		d := day.New(2014, time.February, 7)
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

	Convey("NewTime", t, func() {
		tm := time.Now().AddDate(-10, -5, -1)
		d := day.NewTime(tm)
		So(tm.Year(), ShouldEqual, d.Year())
		So(tm.Month(), ShouldEqual, d.Month())
		So(tm.Day(), ShouldEqual, d.Day())
	})

	Convey("Now", t, func() {
		tm := time.Now()
		d := day.Now()
		So(tm.Year(), ShouldEqual, d.Year())
		So(tm.Month(), ShouldEqual, d.Month())
		So(tm.Day(), ShouldEqual, d.Day())
	})

	Convey("Add", t, func() {
		d := day.New(2016, 2, 29)
		d = d.Add(24 * time.Hour)
		So(d.Year(), ShouldEqual, 2016)
		So(d.Month(), ShouldEqual, time.March)
		So(d.Day(), ShouldEqual, 1)
	})

	Convey("AddDate", t, func() {
		d := day.New(2016, 2, 29)
		d = d.AddDate(0, 0, 1)
		So(d.Year(), ShouldEqual, 2016)
		So(d.Month(), ShouldEqual, time.March)
		So(d.Day(), ShouldEqual, 1)
	})

	Convey("AddDate2", t, func() {
		d := day.New(2016, 2, 29)
		d = d.AddDate2(0, 1, 0)
		So(d.Year(), ShouldEqual, 2016)
		So(d.Month(), ShouldEqual, time.March)
		So(d.Day(), ShouldEqual, 31)
	})

	Convey("After", t, func() {
		d := day.New(2014, time.February, 28)
		before := day.New(2014, time.February, 27)
		after := day.New(2015, time.January, 11)

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
		d := day.New(2015, time.November, 1)
		So(d.Contains(tm, loc), ShouldBeTrue)
	})

	Convey("Equal", t, func() {
		d := day.New(2015, time.November, 7)
		e := day.New(2015, time.November, 8)
		So(d.Equal(e), ShouldBeFalse)
		So(d.Equal(d), ShouldBeTrue)
	})

	Convey("String", t, func() {
		m := day.New(2015, time.November, 17)
		So(m.String(), ShouldEqual, "2015-11-17")
	})
}
