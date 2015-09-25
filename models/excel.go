package models

import (
	"fmt"
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

// Adds a row to xlsx sheet by consuming a pointer to
// InvoiceActivation model based store.
func AddRowXlsx(sheet *xlsx.Sheet, invActivation *InvoiceActivation) error {
	row := sheet.AddRow()
	row.AddCell()

	cell := row.AddCell()
	cell.Value = invActivation.Machine.Name

	// TODO: Implement FastBill product ID
	cell = row.AddCell()
	cell.Value = "Undefined"

	cell = row.AddCell()
	if invActivation.TimeStart.Unix() > 0 {
		cell.Value = invActivation.TimeStart.Format(time.RFC1123)
	}

	cell = row.AddCell()
	cell.SetFloatWithFormat(invActivation.MachineUsage, FORMAT_4_DIGIT)

	cell = row.AddCell()
	cell.Value = invActivation.Machine.PriceUnit

	cell = row.AddCell()
	cell.SetFloatWithFormat(float64(invActivation.Machine.Price), FORMAT_2_DIGIT)

	cell = row.AddCell()
	cell.SetFloatWithFormat(invActivation.PriceTotalExclDisc(), FORMAT_2_DIGIT)

	cell = row.AddCell()
	cell.Value = invActivation.MembershipStr()

	cell = row.AddCell()
	if priceTotalDisc, err := invActivation.PriceTotalDisc(); err == nil {
		cell.SetFloatWithFormat(priceTotalDisc, FORMAT_2_DIGIT)
		cell.SetStyle(boldStyle())
		return nil
	} else {
		return fmt.Errorf("PriceTotalDisc: %v", err)
	}
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
func createXlsxFile(filePath string, invoice *Invoice,
	invSummarry *InvoiceSummary) error {

	sort.Sort(invSummarry)
	userSummaries := &(*invSummarry).UserSummaries

	// Create a xlsx file if there
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet = file.AddSheet("User Summaries")

	// Create header
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "Fab Lab Machine Usage Summary"

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "Period Start Date"
	cell = row.AddCell()
	cell.Value = invSummarry.PeriodStartTime.Format("2006-01-02")

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "Period End Date"
	cell = row.AddCell()
	cell.Value = invSummarry.PeriodEndTime.Format("2006-01-02")

	row = sheet.AddRow()
	row = sheet.AddRow()

	// Fill the xlsx sheet
	for usrSumIter := 0; usrSumIter < len(*userSummaries); usrSumIter++ {
		userSummary := (*userSummaries)[usrSumIter]

		memberships, err := GetUserMemberships(userSummary.User.Id)
		if err != nil {
			return fmt.Errorf("GetUserMemberships: %v", err)
		}

		if len(userSummary.Activations) == 0 && (memberships == nil || len(memberships.Data) == 0) {
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
		cell.Value = strconv.Itoa(userSummary.User.ClientId)

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
		activations := InvoiceActivationsXlsx(userSummary.Activations)
		sort.Stable(activations)
		for _, activation := range activations {
			sumTotal += activation.PriceTotalExclDisc()
			if priceTotalDisc, err := activation.PriceTotalDisc(); err == nil {
				sumTotalDisc += priceTotalDisc
			} else {
				return err
			}

		}

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Activations By Machine"
		AddRowActivationsHeaderXlsx(sheet)

		if summarizedByMachine, err := userSummary.Activations.SummarizedByMachine(); err == nil {
			for _, summed := range summarizedByMachine {
				if err := AddRowXlsx(sheet, summed); err != nil {
					return fmt.Errorf("AddRowXlsx: %v", err)
				}
			}
		} else {
			return fmt.Errorf("SummarizedByMachine: %v", err)
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

		for _, activation := range activations {
			if err := AddRowXlsx(sheet, activation); err != nil {
				return fmt.Errorf("AddRowXlsx: %v", err)
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
