package products

import (
	"github.com/kr15h/fabsmith/models"
	"github.com/kr15h/fabsmith/models/products"
	"github.com/kr15h/fabsmith/tests/assert"
	"github.com/kr15h/fabsmith/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func init() {
	setup.ConfigDB()
}

func TestTutorProducts(t *testing.T) {
	Convey("Testing Tutor product model", t, func() {

		Reset(setup.ResetDB)

		Convey("CreateTutor and GetTutor", func() {
			t := &products.Tutor{
				Product: products.Product{
					UserId: 123,
				},
			}
			t, err1 := products.CreateTutor(t)
			t, err2 := products.GetTutor(t.Product.Id)
			assert.NoErrors(err1, err2)
			So(t.Product.UserId, ShouldEqual, 123)
		})

		Convey("GetAllTutors", func() {
			t1, err1 := products.CreateTutor(&products.Tutor{})
			t2, err2 := products.CreateTutor(&products.Tutor{})
			ts, err := products.GetAllTutors()
			assert.NoErrors(err1, err2, err)
			So(len(ts), ShouldEqual, 2)
			So(t1.Product.Id, ShouldEqual, ts[0].Product.Id)
			So(t2.Product.Id, ShouldEqual, ts[1].Product.Id)
		})

		Convey("Update", func() {
			u := models.User{
				FirstName: "Roland",
				LastName:  "Kaiser",
				Email:     "roland.kaiser@signal-iduna.de",
			}
			uid, err := models.CreateUser(&u)
			t := &products.Tutor{
				Product: products.Product{
					UserId: uid,
				},
			}
			t, err = products.CreateTutor(t)

			u2 := models.User{
				FirstName: "Peter",
				LastName:  "Lustig",
				Email:     "peter.lustig@wdr.de",
			}
			t.Product.UserId, err = models.CreateUser(&u2)
			err2 := t.Update()

			t, err3 := products.GetTutor(t.Product.Id)
			assert.NoErrors(err, err2, err3)
			So(t.Product.UserId, ShouldEqual, u2.Id)
		})

	})
}
