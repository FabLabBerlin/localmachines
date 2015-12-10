package products

import (
	"github.com/kr15h/fabsmith/models/products"
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
			s, err := products.CreateSpace("foo")
			s2, err2 := products.GetSpace(s.Product.Id)
			assert.NoErrors(err, err2)
			So(s2.Product.Name, ShouldEqual, "foo")
		})

		Convey("GetAllSpaces", func() {
			s1, err1 := products.CreateSpace("foo")
			s2, err2 := products.CreateSpace("bar")
			l, err := products.GetAllSpaces()
			assert.NoErrors(err1, err2, err)
			So(len(l), ShouldEqual, 2)
			So(s1.Product.Id, ShouldEqual, l[0].Product.Id)
			So(s2.Product.Id, ShouldEqual, l[1].Product.Id)
		})

		Convey("Update", func() {
			s, err := products.CreateSpace("foobar")
			s.Product.Name = "foobaz"
			err2 := s.Update()
			s2, err3 := products.GetSpace(s.Product.Id)
			assert.NoErrors(err, err2, err3)
			So("foobaz", ShouldEqual, s2.Product.Name)
		})
	})
}
