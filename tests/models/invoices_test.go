package modelTest

import (
	"os"
	"sort"
	"testing"
	"time"

	"github.com/kr15h/fabsmith/models"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tealeg/xlsx"
)

var TIME_START = time.Now()

func init() {
	ConfigDB()
}

func CreateTestPurchase(machineId int64, machineName string,
	minutes time.Duration, pricePerMinute float64) *models.Purchase {

	machine := models.Machine{}
	machine.Id = machineId
	machine.Name = machineName
	machine.PriceUnit = "minute"
	machine.Price = pricePerMinute

	invAct := &models.Purchase{
		Activation: &models.Activation{
			TimeStart:                   TIME_START,
			TimeEnd:                     TIME_START.Add(time.Minute * time.Duration(minutes)),
			CurrentMachinePrice:         pricePerMinute,
			CurrentMachinePriceUnit:     "minute",
			CurrentMachinePriceCurrency: "€",
		},
		Machine:      &machine,
		MachineUsage: minutes,
		Memberships: []*models.Membership{
			&models.Membership{
				Id:                    42,
				Title:                 "Half price",
				ShortName:             "HP",
				MachinePriceDeduction: 50,
				AffectedMachines:      "[22]",
			},
		},
	}
	invAct.TotalPrice = models.PriceTotalExclDisc(invAct)
	var err error
	invAct.DiscountedTotal, err = models.PriceTotalDisc(invAct)
	if err != nil {
		panic(err.Error())
	}
	return invAct
}

func TestInvoiceActivation(t *testing.T) {
	Convey("Testing InvoiceActivation model", t, func() {
		Reset(ResetDB)
		Convey("Testing MembershipStr", func() {
			invAct := CreateTestPurchase(22, "Lasercutter",
				time.Duration(12)*time.Minute, 0.5)
			So(invAct.MembershipStr(), ShouldEqual, "HP (50%)")
		})
		Convey("Testing PriceTotalExclDisc", func() {
			invAct := CreateTestPurchase(22, "Lasercutter",
				time.Duration(12)*time.Minute, 0.5)
			So(models.PriceTotalExclDisc(invAct), ShouldEqual, 6)
		})
		Convey("Testing PriceTotalDisc", func() {
			invAct := CreateTestPurchase(22, "Lasercutter",
				time.Duration(12)*time.Minute, 0.5)
			if priceTotalDisc, err := models.PriceTotalDisc(invAct); err == nil {
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
			models.AddRowActivationsHeaderXlsx(sheet)
			models.AddRowXlsx(sheet, invAct)
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
		Reset(ResetDB)

		Convey("Testing SummarizedByMachine", func() {
			invs := models.Purchases{
				Data: []*models.Purchase{
					CreateTestPurchase(22, "Lasercutter", time.Duration(12)*time.Minute, 0.5),
					CreateTestPurchase(22, "Lasercutter", time.Duration(13)*time.Minute, 0.25),
					CreateTestPurchase(23, "CNC Router", time.Duration(12)*time.Minute, 0.8),
					CreateTestPurchase(23, "CNC Router", time.Duration(12)*time.Minute, 0.8),
					CreateTestPurchase(23, "CNC Router", time.Duration(12)*time.Minute, 0.8),
				},
			}

			if invs, err := invs.SummarizedByMachine(); err == nil {
				sort.Stable(invs)
				So(len(invs.Data), ShouldEqual, 2)
				So(invs.Data[0].Machine.Name, ShouldEqual, "CNC Router")
				So(invs.Data[0].TotalPrice, ShouldAlmostEqual, 36*0.8, 0.000001)
				So(invs.Data[0].TotalPrice, ShouldAlmostEqual, models.PriceTotalExclDisc(invs.Data[0]), 0.000001)
				So(invs.Data[0].DiscountedTotal, ShouldAlmostEqual, 36*0.8, 0.000001)
				So(invs.Data[1].DiscountedTotal, ShouldAlmostEqual, 0.5*(12*0.5+13*0.25), 0.000001)
				if priceTotalDisc, err := models.PriceTotalDisc(invs.Data[0]); err == nil {
					So(invs.Data[0].DiscountedTotal, ShouldAlmostEqual, priceTotalDisc, 0.000001)
				} else {
					panic(err.Error())
				}
				So(invs.Data[1].Machine.Name, ShouldEqual, "Lasercutter")
			} else {
				panic(err.Error())
			}
		})
	})

	Convey("When creating invoice with CreateInvoice", t, func() {

		endTime := time.Now()
		startTime := endTime.AddDate(0, -1, 0)

		invoice, err := models.CreateInvoice(startTime, endTime)

		Convey("It should not cause any error", func() {
			So(err, ShouldBeNil)
		})

		Convey("The returned pointer to Invoice store should be valid", func() {
			So(invoice, ShouldNotBeNil)
		})

		Convey("When reading back the created invoice from the db", func() {
			var readbackInvoice *models.Invoice
			readbackInvoice, err = models.GetInvoice(invoice.Id)

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
			invoices, err := models.GetAllInvoices()

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
