package modelTest

import (
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

func CreateTestInvActivation(machineName string, minutes, pricePerMinute float64) *models.InvoiceActivation {
	invAct := &models.InvoiceActivation{
		TimeStart:           TIME_START,
		TimeEnd:             TIME_START.Add(time.Minute * time.Duration(minutes)),
		MachineName:         machineName,
		MachineProductId:    "Undefined",
		MachineUsage:        minutes,
		MachineUsageUnit:    "minute",
		MachinePricePerUnit: pricePerMinute,
		Memberships: []*models.Membership{
			&models.Membership{
				Id:                    42,
				Title:                 "Half price",
				ShortName:             "HP",
				MachinePriceDeduction: 50,
			},
		},
	}
	return invAct
}

func TestInvoiceActivation(t *testing.T) {
	Convey("Testing InvoiceActivation model", t, func() {
		Reset(ResetDB)
		Convey("Testing MembershipStr", func() {
			invAct := CreateTestInvActivation("Lasercutter", 12, 0.5)
			So(invAct.MembershipStr(), ShouldEqual, "HP (50%)")
		})
		Convey("Testing PriceTotalExclDisc", func() {
			invAct := CreateTestInvActivation("Lasercutter", 12, 0.5)
			So(invAct.PriceTotalExclDisc(), ShouldEqual, 6)
		})
		Convey("Testing PriceTotalDisc", func() {
			invAct := CreateTestInvActivation("Lasercutter", 12, 0.5)
			So(invAct.PriceTotalDisc(), ShouldEqual, 3)
		})
		Convey("Testing AddRowXlsx", func() {
			inv := models.Invoice{}
			testTable := [][]interface{}{
				[]interface{}{"", "Machine Name", "Product ID", "Start Time", "Usage", "Usage Unit", "€ per Unit", "Total €", "Memberships", "Discounted €"},
				[]interface{}{"", "Lasercutter", "Undefined", TIME_START.Format(time.RFC1123), "12", "minute", "0.50", "6.00", "HP (50%)", "3.00"},
			}
			invAct := CreateTestInvActivation("Lasercutter", 12, 0.5)
			file := xlsx.NewFile()
			sheet := file.AddSheet("User Summaries")
			inv.AddRowActivationsHeaderXlsx(sheet)
			invAct.AddRowXlsx(sheet)
			m := 2
			So(len(sheet.Rows), ShouldEqual, m)
			for i := 0; i < m; i++ {
				cells := sheet.Rows[i].Cells
				n := 10
				So(len(cells), ShouldEqual, n)
				for j := 0; j < n; j++ {
					So(cells[j].String(), ShouldEqual, testTable[i][j])
				}
			}
		})
	})

	Convey("Testing InvoiceActivations model", t, func() {
		Reset(ResetDB)
		Convey("Testing SummarizedByMachine", func() {
			invs := models.InvoiceActivations{
				CreateTestInvActivation("Lasercutter", 12, 0.5),
				CreateTestInvActivation("Lasercutter", 13, 0.25),
				CreateTestInvActivation("CNC Router", 12, 0.8),
				CreateTestInvActivation("CNC Router", 12, 0.8),
				CreateTestInvActivation("CNC Router", 12, 0.8),
			}
			invs = invs.SummarizedByMachine()
			sort.Stable(invs)
			So(len(invs), ShouldEqual, 2)
			So(invs[0].MachineName, ShouldEqual, "CNC Router")
			So(invs[0].TotalPrice, ShouldAlmostEqual, 36*0.8, 0.000001)
			So(invs[0].TotalPrice, ShouldAlmostEqual, invs[0].PriceTotalExclDisc(), 0.000001)
			So(invs[0].DiscountedTotal, ShouldAlmostEqual, 0.5*36*0.8, 0.000001)
			So(invs[0].DiscountedTotal, ShouldAlmostEqual, invs[0].PriceTotalDisc(), 0.000001)
			So(invs[1].MachineName, ShouldEqual, "Lasercutter")
		})
	})
}
