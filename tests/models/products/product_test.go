package products

import (
	"fmt"
	"github.com/kr15h/fabsmith/models/products"
	"github.com/kr15h/fabsmith/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestProducts(t *testing.T) {
	Convey("Testing Product model", t, func() {

		Reset(setup.ResetDB)

		Convey("Create and Get", func() {
			p := &products.Product{
				Name: "foo",
			}
			id, err := products.Create(p)
			if err != nil {
				panic(err.Error())
			}
			p, err = products.Get(id)
			if err != nil {
				panic(err.Error())
			}
			So(p.Name, ShouldEqual, "foo")
		})

		Convey("GetAll", func() {
			id1, err1 := products.Create(&products.Product{
				Name: "bar",
			})
			id2, err2 := products.Create(&products.Product{
				Name: "baz",
			})
			ps, err := products.GetAll()
			if err1 != nil || err2 != nil || err != nil {
				panic(fmt.Sprintf("Errors: %v, %v, %v", err1, err2, err))
			}
			So(len(ps), ShouldEqual, 2)
			So(id1, ShouldEqual, ps[0].Id)
			So(id2, ShouldEqual, ps[1].Id)
		})

		Convey("Update", func() {
			p := &products.Product{
				Name: "foo",
			}
			id, err := products.Create(p)
			p.Name = "foobar"
			err2 := p.Update()
			p, err3 := products.Get(id)
			if err != nil || err2 != nil || err3 != nil {
				panic(fmt.Sprintf("Errors: %v, %v, %v", err, err2, err3))
			}
			So("foobar", ShouldEqual, p.Name)
		})

		Convey("Archive", func() {
			p := &products.Product{
				Name: "foo",
			}
			_, err := products.Create(p)
			err2 := p.Archive()
			ps, err3 := products.GetAll()
			if err != nil || err2 != nil || err3 != nil {
				panic(fmt.Sprintf("Errors: %v, %v, %v", err, err2, err3))
			}
			So(len(ps), ShouldEqual, 0)
		})

	})
}
