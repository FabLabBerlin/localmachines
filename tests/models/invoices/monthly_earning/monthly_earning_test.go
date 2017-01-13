package invoices

import (
	"os"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/invoices/monthly_earning"
	"github.com/FabLabBerlin/localmachines/models/locations"
	//"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/settings"
	"github.com/FabLabBerlin/localmachines/tests/models/util"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tealeg/xlsx"
)

var TIME_START = util.TIME_START

func init() {
	setup.ConfigDB()
}

func TestInvoiceActivation(t *testing.T) {
	Convey("Testing InvoiceActivation model", t, func() {
		Reset(setup.ResetDB)
		Convey("Testing PriceTotalExclDisc", func() {
			invAct := util.CreateTestPurchase(22, 0, "Lasercutter",
				time.Duration(12)*time.Minute, 0.5)
			So(purchases.PriceTotalExclDisc(invAct), ShouldEqual, 6)
		})
		Convey("Testing AddRowXlsx", func() {
			currency := "$"
			if _, err := settings.Create(&settings.Setting{
				LocationId:  1,
				Name:        settings.CURRENCY,
				ValueString: &currency,
			}); err != nil {
				panic(err.Error())
			}

			invAct := util.CreateTestPurchase(22, 0, "Lasercutter",
				time.Duration(12)*time.Minute, 0.5)
			file := xlsx.NewFile()
			sheet, _ := file.AddSheet("User Summaries")
			monthly_earning.AddRowActivationsHeaderXlsx(&monthly_earning.MonthlyEarning{
				LocationId: 1,
				Currency:   "$",
			}, sheet)
			loc := &locations.Location{
				Id:       1,
				Timezone: "Europe/Berlin",
			}
			monthly_earning.AddRowXlsx(loc, sheet, invAct)
			numRows := 2

			Convey("Number of rows in xlsx sheed should be correct", func() {
				So(len(sheet.Rows), ShouldEqual, numRows)
			})
		})
	})

	Convey("When creating Monthly Earning with monthly_earning.New()", t, func() {
		endTime := time.Now()
		startTime := time.Now()
		interval := lib.Interval{
			MonthFrom: int(startTime.Month()),
			YearFrom:  startTime.Year(),
			MonthTo:   int(endTime.Month()),
			YearTo:    endTime.Year(),
		}

		me, err := monthly_earning.Create(1, interval)

		Convey("It should not cause any error", func() {
			So(err, ShouldBeNil)
		})

		Convey("The returned pointer to object should be valid", func() {
			So(me, ShouldNotBeNil)
		})

		Convey("When reading back the created monthly earning from the db", func() {
			dbMe, err := monthly_earning.Get(me.Id)

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("The pointer to the read invoice should be valid", func() {
				So(dbMe, ShouldNotBeNil)
			})

			Convey("File path should be there", func() {
				So(dbMe.FilePath, ShouldEqual, me.FilePath)
			})
		})

		Convey("When trying to get all monthly earnings", func() {
			mes, err := monthly_earning.GetAllAt(1)

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("The number of returned monthly earnings should be more than 0", func() {
				So(len(mes), ShouldBeGreaterThan, 0)
			})
		})

		// Remove temp files directory used for the monthly earning files
		os.RemoveAll("files")
	})

}
