package tests

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOnePlusOne(t *testing.T) {
	Convey("Testing 1+1", t, func() {
		res := 1 + 1
		Convey("Should equal 2", func() {
			So(res, ShouldEqual, 2)
		})
	})
}
