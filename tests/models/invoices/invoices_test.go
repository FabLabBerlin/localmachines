package invoices

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/invoices"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tealeg/xlsx"
)

var TIME_START = time.Now()

func init() {
	setup.ConfigDB()
}

func CreateTestPurchase(machineId int64, machineName string,
	minutes time.Duration, pricePerMinute float64) *purchases.Purchase {

	machine := models.Machine{}
	machine.Id = machineId
	machine.Name = machineName
	machine.PriceUnit = "minute"
	machine.Price = pricePerMinute

	invAct := &purchases.Purchase{
		Type:         purchases.TYPE_ACTIVATION,
		TimeStart:    TIME_START,
		TimeEnd:      TIME_START.Add(minutes),
		PricePerUnit: pricePerMinute,
		PriceUnit:    "minute",
		Quantity:     minutes.Minutes(),
		Machine:      &machine,
		MachineId:    machineId,
		MachineUsage: minutes,
		Memberships: []*models.Membership{
			&models.Membership{
				Id:                    42,
				Title:                 "Half price",
				ShortName:             "HP",
				MachinePriceDeduction: 50,
				AffectedMachines:      fmt.Sprintf("[%v]", machineId),
			},
		},
	}
	invAct.TotalPrice = purchases.PriceTotalExclDisc(invAct)
	var err error
	invAct.DiscountedTotal, err = purchases.PriceTotalDisc(invAct)
	if err != nil {
		panic(err.Error())
	}
	return invAct
}

func TestInvoiceActivation(t *testing.T) {
	Convey("Testing InvoiceActivation model", t, func() {
		Reset(setup.ResetDB)
		Convey("Testing MembershipStr", func() {
			invAct := CreateTestPurchase(22, "Lasercutter",
				time.Duration(12)*time.Minute, 0.5)
			So(invAct.MembershipStr(), ShouldEqual, "HP (50%)")
		})
		Convey("Testing PriceTotalExclDisc", func() {
			invAct := CreateTestPurchase(22, "Lasercutter",
				time.Duration(12)*time.Minute, 0.5)
			So(purchases.PriceTotalExclDisc(invAct), ShouldEqual, 6)
		})
		Convey("Testing PriceTotalDisc", func() {
			invAct := CreateTestPurchase(22, "Lasercutter",
				time.Duration(12)*time.Minute, 0.5)
			if priceTotalDisc, err := purchases.PriceTotalDisc(invAct); err == nil {
				So(priceTotalDisc, ShouldEqual, 3)
			} else {
				panic(err.Error())
			}
		})
		Convey("Testing AddRowXlsx", func() {
			testTable := [][]interface{}{
				[]interface{}{"", "Machine Name", "Product ID",
					"Start Time", "Usage", "Usage Unit", "€ per Unit",
					"Total €", "Memberships", "Discounted €"},

				[]interface{}{"", "Lasercutter", "Undefined",
					TIME_START.Format(time.RFC1123), "12", "minute",
					"0.50", "6.00", "HP (50%)", "3.00"},
			}

			invAct := CreateTestPurchase(22, "Lasercutter",
				time.Duration(12)*time.Minute, 0.5)
			file := xlsx.NewFile()
			sheet, _ := file.AddSheet("User Summaries")
			invoices.AddRowActivationsHeaderXlsx(sheet)
			invoices.AddRowXlsx(sheet, invAct)
			numRows := 2

			Convey("Number of rows in xlsx sheed should be correct", func() {
				So(len(sheet.Rows), ShouldEqual, numRows)
			})

			Convey("The number and content of cells should be correct", func() {
				for i := 0; i < numRows; i++ {
					cells := sheet.Rows[i].Cells
					numCells := 10

					So(len(cells), ShouldEqual, numCells)

					for j := 0; j < numCells; j++ {
						So(cells[j].String(), ShouldEqual, testTable[i][j])
					}
				}
			})

		})
	})

	Convey("Testing InvoiceActivations model", t, func() {
		Reset(setup.ResetDB)
		// TODO: Write tests for this
	})

	Convey("When creating invoice with CreateInvoice", t, func() {

		endTime := time.Now()
		startTime := endTime.AddDate(0, -1, 0)

		invoice, err := invoices.Create(startTime, endTime)

		Convey("It should not cause any error", func() {
			So(err, ShouldBeNil)
		})

		Convey("The returned pointer to Invoice store should be valid", func() {
			So(invoice, ShouldNotBeNil)
		})

		Convey("When reading back the created invoice from the db", func() {
			var readbackInvoice *invoices.Invoice
			readbackInvoice, err = invoices.Get(invoice.Id)

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("The pointer to the read invoice should be valid", func() {
				So(readbackInvoice, ShouldNotBeNil)
			})

			Convey("The read back invoice start and end time should be correct", func() {
				So(readbackInvoice.PeriodFrom, ShouldHappenWithin,
					time.Duration(1)*time.Second, startTime)
				So(readbackInvoice.PeriodTo, ShouldHappenWithin,
					time.Duration(1)*time.Second, endTime)
			})

			Convey("File path should be there", func() {
				So(readbackInvoice.FilePath, ShouldEqual, invoice.FilePath)
			})
		})

		Convey("When trying to get all invoices", func() {
			invoices, err := invoices.GetAll()

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("The number of returned invoices should be more than 0", func() {
				So(len(*invoices), ShouldBeGreaterThan, 0)
			})
		})

		// Remove temp files directory used for the invoice files
		os.RemoveAll("files")
	})

}
