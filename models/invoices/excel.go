package invoices

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/tealeg/xlsx"
	"sort"
	"strconv"
	"time"
)

const (
	FORMAT_2_DIGIT = "#,##0.00"
	FORMAT_4_DIGIT = "#,####0.0000"

	BLUE   = "FF63C5E5"
	RED    = "FFEA535D"
	YELLOW = "FFFBE16B"
	GREEN  = "FF92D050"
)

type PurchasesXlsx []*purchases.Purchase

func (this PurchasesXlsx) Len() int {
	return len(this)
}

func (this PurchasesXlsx) Less(i, j int) bool {
	timeStartI := (*this[i]).TimeStart
	timeStartJ := (*this[j]).TimeStart
	anyNil := (*this[i]).Machine == nil || (*this[j]).Machine == nil
	if !anyNil && (*this[i]).Machine.Name < (*this[j]).Machine.Name {
		return true
	} else if !anyNil && (*this[j]).Machine.Name < (*this[i]).Machine.Name {
		return false
	} else {
		return timeStartI.Before(timeStartJ)
	}
}

func (this PurchasesXlsx) Swap(i, j int) {
	*this[i], *this[j] = *this[j], *this[i]
}

// Adds a row to xlsx sheet by consuming a pointer to
// InvoiceActivation model based store.
func AddRowXlsx(sheet *xlsx.Sheet, purchase *purchases.Purchase) (err error) {
	timeStart := purchase.TimeStart
	totalPrice := purchase.TotalPrice
	discountedTotal := purchase.DiscountedTotal

	row := sheet.AddRow()
	row.AddCell()

	cell := row.AddCell()
	cell.Value = purchase.ProductName()

	// TODO: Implement FastBill product ID
	cell = row.AddCell()
	cell.Value = "Undefined"

	cell = row.AddCell()
	if timeStart.Unix() > 0 {
		cell.Value = timeStart.Format(time.RFC1123)
	}

	cell = row.AddCell()
	cell.SetFloatWithFormat(purchase.Quantity, FORMAT_4_DIGIT)

	cell = row.AddCell()
	cell.Value = purchase.PriceUnit

	cell = row.AddCell()
	cell.SetFloatWithFormat(purchase.PricePerUnit, FORMAT_2_DIGIT)

	cell = row.AddCell()
	cell.SetFloatWithFormat(totalPrice, FORMAT_2_DIGIT)

	cell = row.AddCell()
	if cell.Value, err = purchase.MembershipStr(); err != nil {
		return fmt.Errorf("membership string: %v", err)
	}

	cell = row.AddCell()
	cell.SetFloatWithFormat(discountedTotal, FORMAT_2_DIGIT)
	cell.SetStyle(boldStyle())
	return nil
}

// Adds header row to existing xlsx sheet.
func AddRowActivationsHeaderXlsx(sheet *xlsx.Sheet) {
	row := sheet.AddRow()
	row.AddCell()
	cell := row.AddCell()
	cell.Value = "Machine Name"
	cell = row.AddCell()
	cell.Value = "Product ID"
	cell = row.AddCell()
	cell.Value = "Start Time"
	cell = row.AddCell()
	cell.Value = "Usage"
	cell = row.AddCell()
	cell.Value = "Usage Unit"
	cell = row.AddCell()
	cell.Value = "€ per Unit"
	cell = row.AddCell()
	cell.Value = "Total €"
	cell = row.AddCell()
	cell.Value = "Memberships"
	cell = row.AddCell()
	cell.Value = "Discounted €"
}

// Adds an empty row.
func addSeparationRowXlsx(sheet *xlsx.Sheet) {
	row := sheet.AddRow()
	style := colorStyle(YELLOW)
	for i := 0; i < 11; i++ {
		cell := row.AddCell()
		cell.SetStyle(style)
	}
}

// Creates a xlsx file.
func createXlsxFile(filePath string, invoice *Invoice) error {
	sort.Sort(invoice)
	userSummaries := invoice.UserSummaries

	// Create a xlsx file if there
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, _ = file.AddSheet("User Summaries")

	// Create header
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "Fab Lab Machine Usage Summary"

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "Period Start Date"
	cell = row.AddCell()
	cell.Value = invoice.PeriodFrom.Format("2006-01-02")

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "Period End Date"
	cell = row.AddCell()
	cell.Value = invoice.PeriodTo.Format("2006-01-02")

	row = sheet.AddRow()
	row = sheet.AddRow()

	// Fill the xlsx sheet
	for _, userSummary := range userSummaries {

		memberships, err := models.GetUserMemberships(userSummary.User.Id)
		if err != nil {
			return fmt.Errorf("GetUserMemberships: %v", err)
		}

		if len(userSummary.Purchases.Data) == 0 &&
			(memberships == nil || len(memberships.Data) == 0) {
			// nothing to bill
			continue
		}

		addSeparationRowXlsx(sheet)
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "User"
		cell = row.AddCell()
		cell.Value = userSummary.User.FirstName
		cell = row.AddCell()
		cell.Value = userSummary.User.LastName
		cell = row.AddCell()
		cell.Value = userSummary.User.Email

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Fastbill User Id"
		cell = row.AddCell()
		cell.SetStyle(colorStyle(RED))
		cell.Value = strconv.FormatInt(userSummary.User.ClientId, 10)

		// User Billing Address
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Billing Address"
		cell = row.AddCell()
		cell.Value = userSummary.User.InvoiceAddr

		// User Zip Code
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Zip Code"
		cell = row.AddCell()
		cell.Value = userSummary.User.ZipCode

		// User City
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "City"
		cell = row.AddCell()
		cell.Value = userSummary.User.City

		// User Country
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Country"
		cell = row.AddCell()
		cell.Value = userSummary.User.CountryCode

		// User Phone
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Phone"
		cell = row.AddCell()
		cell.Value = userSummary.User.Phone

		// Company
		if userSummary.User.Company != "" {
			row = sheet.AddRow()
			cell = row.AddCell()
			cell.Value = "Company"
			cell = row.AddCell()
			cell.Value = userSummary.User.Company
		}

		// User Comments
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Comments"
		cell = row.AddCell()
		cell.Value = userSummary.User.Comments

		if memberships != nil {
			sheet.AddRow()
			sheet.AddRow()
			row = sheet.AddRow()
			cell = row.AddCell()
			cell.Value = "Memberships"
			row = sheet.AddRow()
			row.AddCell()
			cell = row.AddCell()
			cell.SetStyle(boldStyle())
			cell.Value = "Title"
			cell = row.AddCell()
			cell.Value = "Start Date"
			cell = row.AddCell()
			cell.Value = "End Date"
			cell = row.AddCell()
			cell.SetStyle(boldStyle())
			cell.Value = "Monthly Price / €"
			cell = row.AddCell()
			cell.Value = "Duration Unit"
			cell = row.AddCell()
			cell.Value = "Machine Price Deduction"
			for _, m := range memberships.Data {
				row = sheet.AddRow()
				row.AddCell()
				cell = row.AddCell()
				cell.SetStyle(colorStyle(BLUE))
				cell.Value = m.Title
				cell = row.AddCell()
				cell.Value = m.StartDate.Format(time.RFC1123)
				cell = row.AddCell()
				cell.Value = m.EndDate.Format(time.RFC1123)
				cell = row.AddCell()
				cell.SetFloatWithFormat(float64(m.MonthlyPrice), FORMAT_2_DIGIT)
				cell.SetStyle(colorStyle(GREEN))
				cell = row.AddCell()
				cell.Value = m.Unit
				cell = row.AddCell()
				cell.Value = strconv.Itoa(m.MachinePriceDeduction) + "%"
			}
			sheet.AddRow()
			sheet.AddRow()
		}

		sumTotal := 0.0
		sumTotalDisc := 0.0
		ps := PurchasesXlsx(userSummary.Purchases.Data)
		sort.Stable(ps)
		for _, purchase := range ps {
			sumTotal += purchase.TotalPrice
			sumTotalDisc += purchase.DiscountedTotal

		}

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Activations By Machine"
		AddRowActivationsHeaderXlsx(sheet)

		byProductNameAndPricePerUnit := userSummary.byProductNameAndPricePerUnit()

		for productName, byPricePerUnit := range byProductNameAndPricePerUnit {
			for pricePerUnit, ps := range byPricePerUnit {
				var usage float64
				var usageUnit string
				var totalPriceExclDisc float64
				var discPrice float64
				var membershipStr string
				for _, purchase := range ps {
					usageUnit = purchase.PriceUnit
					usage += purchase.Quantity
					totalPriceExclDisc += purchases.PriceTotalExclDisc(purchase)
					priceDisc, err := purchases.PriceTotalDisc(purchase)
					if err != nil {
						return fmt.Errorf("PriceTotalDisc: %v", err)
					}
					discPrice += priceDisc
					if membershipStr, err = purchase.MembershipStr(); err != nil {
						return fmt.Errorf("membership string: %v", err)
					}
				}
				row = sheet.AddRow()
				row.AddCell()
				row.AddCell().Value = productName
				row.AddCell().Value = "Undefined"
				row.AddCell()
				row.AddCell().SetFloatWithFormat(usage, FORMAT_4_DIGIT)
				row.AddCell().Value = usageUnit
				row.AddCell().SetFloatWithFormat(pricePerUnit, FORMAT_2_DIGIT)
				row.AddCell().SetFloatWithFormat(totalPriceExclDisc, FORMAT_2_DIGIT)
				row.AddCell().Value = membershipStr
				row.AddCell().SetFloatWithFormat(discPrice, FORMAT_2_DIGIT)
			}
		}

		printTotal := func(totalColor string) {
			row = sheet.AddRow()
			row.AddCell()
			row.AddCell()
			row.AddCell()
			row.AddCell()
			row.AddCell()
			row.AddCell()
			cell = row.AddCell()
			cell.Value = "Subtotal €"
			cell = row.AddCell()
			cell.SetFloatWithFormat(sumTotal, FORMAT_2_DIGIT)
			cell = row.AddCell()
			cell.SetStyle(boldStyle())
			cell.Value = "Discounted €"
			cell = row.AddCell()
			cell.SetFloatWithFormat(sumTotalDisc, FORMAT_2_DIGIT)
			cell.SetStyle(colorStyle(totalColor))
		}
		printTotal(BLUE)

		sheet.AddRow()
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Activations"
		AddRowActivationsHeaderXlsx(sheet)

		for _, purchase := range ps {
			if !purchase.Cancelled {
				if err := AddRowXlsx(sheet, purchase); err != nil {
					return fmt.Errorf("AddRowXlsx: %v", err)
				}
			}
		}
		printTotal(GREEN)

		sheet.AddRow()
	} // for userSummaries

	return file.Save(filePath)
}

// Returns xlsx bold style.
func boldStyle() *xlsx.Style {
	font := xlsx.DefaultFont()
	font.Bold = true
	style := xlsx.NewStyle()
	style.Font = *font
	return style
}

// Returns xlsx colored style.
func colorStyle(color string) *xlsx.Style {
	font := xlsx.DefaultFont()
	font.Bold = true
	style := xlsx.NewStyle()
	style.Fill = *xlsx.NewFill("solid", color, "FF00FF00")
	style.Font = *font
	return style
}
