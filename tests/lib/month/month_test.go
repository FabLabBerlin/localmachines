package month

import (
	"github.com/FabLabBerlin/localmachines/lib/month"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestLibMonth(t *testing.T) {
	Convey("Construct and Copy", t, func() {
		m := month.New(2014, time.February)
		var cpy month.Month
		cpy = m
		So(cpy.Month(), ShouldEqual, time.February)
		So(cpy.Year(), ShouldEqual, 2014)
	})

	Convey("NewString", t, func() {
		for _, s := range []string{"2013-07", "2013-07-01"} {
			m, err := month.NewString(s)
			So(err, ShouldBeNil)
			So(m.Month(), ShouldEqual, time.July)
			So(m.Year(), ShouldEqual, 2013)
		}
	})

	Convey("Add", t, func() {
		m := month.New(2014, time.February)
		So(m.Add(1).String(), ShouldEqual, "2014-03")
		So(m.Add(13).String(), ShouldEqual, "2015-03")
		So(m.Add(26).String(), ShouldEqual, "2016-04")

		m = month.New(2015, time.November)
		So(m.Add(1).String(), ShouldEqual, "2015-12")

		m = month.New(2015, time.December)
		So(m.Add(1).String(), ShouldEqual, "2016-01")
	})

	Convey("After", t, func() {
		m := month.New(2014, time.February)
		before := month.New(2014, time.January)
		after := month.New(2015, time.January)

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
		m := month.New(2015, time.November)
		So(m.Contains(tm, loc), ShouldBeTrue)
	})

	Convey("Equal", t, func() {
		m := month.New(2015, time.November)
		n := month.New(2015, time.December)
		So(m.Equal(n), ShouldBeFalse)
		So(m.Equal(m), ShouldBeTrue)
	})

	Convey("IsZero", t, func() {
		m := month.Month{}
		So(m.IsZero(), ShouldBeTrue)
		m = month.New(1999, time.December)
		So(m.IsZero(), ShouldBeFalse)
	})

	Convey("String", t, func() {
		m := month.New(2015, time.November)
		So(m.String(), ShouldEqual, "2015-11")
	})
}
