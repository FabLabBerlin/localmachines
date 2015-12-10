package purchases

import (
	"github.com/kr15h/fabsmith/models/purchases"
	"github.com/kr15h/fabsmith/tests/assert"
	"github.com/kr15h/fabsmith/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestSpacePurchases(t *testing.T) {
	Convey("Testing Space purchase model", t, func() {

		Reset(setup.ResetDB)

		Convey("CreateSpace and GetSpace", func() {
			s := purchases.NewSpace()
			s.UserId = 234
			err1 := s.Save()
			s, err2 := purchases.GetSpace(s.Id)
			assert.NoErrors(err1, err2)
			So(s.Purchase.UserId, ShouldEqual, 234)
		})

		Convey("GetAllSpaces", func() {
			s1 := purchases.NewSpace()
			err1 := s1.Save()
			s2 := purchases.NewSpace()
			err2 := s2.Save()
			l, err := purchases.GetAllSpace()
			assert.NoErrors(err1, err2, err)
			So(len(l), ShouldEqual, 2)
			So(s1.Id, ShouldEqual, l[0].Id)
			So(s2.Id, ShouldEqual, l[1].Id)
		})

		Convey("Update", func() {
			s := purchases.NewSpace()
			s.UserId = 234
			err := s.Save()
			s.UserId = 456
			err2 := s.Update()
			s, err3 := purchases.GetSpace(s.Id)
			assert.NoErrors(err, err2, err3)
			So(456, ShouldEqual, s.Purchase.UserId)
		})

	})
}
