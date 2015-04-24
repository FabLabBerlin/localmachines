package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/tealeg/xlsx"
	"math/rand"
	"strconv"
	"time"
)

// Invoice entry that is saved in the database.
// Activations field contains a JSON array with activation IDs.
// XlsFile field contains URL to the generated XLSX file.
type Invoice struct {
	Id          int64  `orm:"auto";"pk"`
	Activations string `orm:type(text)`
	FilePath    string `orm:size(255)`
	Created     time.Time
	PeriodFrom  time.Time
	PeriodTo    time.Time
}

func init() {
	orm.RegisterModel(new(Invoice))
}

func (this *Invoice) TableName() string {
	return "invoices"
}

// This is a user activation row that appears in the XLSX file
type InvoiceActivation struct {
	MachineId           int64 // Machine info
	MachineName         string
	MachineProductId    string
	MachineUsage        float64
	MachineUsageUnit    string
	MachinePricePerUnit float64
	UserId              int64 // User info
	UserFullName        string
	Username            string
	UserEmail           string
	UserDebitorNumber   string
	Memberships         []*Membership
	TotalPrice          float64
	DiscountedTotal     float64
}

type UserSummary struct {
	UserId        int64
	UserFullName  string
	Username      string
	UserEmail     string
	DebitorNumber string
	Activations   []*InvoiceActivation
}

type InvoiceSummary struct {
	InvoiceName     string
	PeriodStartTime time.Time
	PeriodEndTime   time.Time
	UserSummaries   []*UserSummary
}

func CreateInvoice(startTime, endTime time.Time) (*Invoice, error) {

	var err error

	// Create invoice object that will be returned
	invoice := Invoice{}

	// Get all uninvoiced activations in the time range
	var activations *[]Activation
	activations, err = invoice.getActivations(startTime, endTime)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get activations: %v", err))
	}

	beego.Trace("Activations str:", invoice.Activations)

	// Enhance activations with user and membership data
	var enhancedActivations *[]*InvoiceActivation
	enhancedActivations, err = invoice.getEnhancedActivations(activations)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get enhanced activations: %v", err))
	}

	// Create user summaries from invoice activations
	var userSummaries *[]*UserSummary
	userSummaries = invoice.getUserSummaries(enhancedActivations)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get user summaries: %v", err))
	}

	for i := 0; i < len(*userSummaries); i++ {
		beego.Trace((*userSummaries)[i].Activations)
	}

	// Create invoice summary
	invSummary := InvoiceSummary{}
	invSummary.InvoiceName = "Invoice"
	invSummary.PeriodStartTime = startTime
	invSummary.PeriodEndTime = endTime
	invSummary.UserSummaries = *userSummaries

	// Create *.xlsx file.
	fileName := invoice.getInvoiceFileName(startTime, endTime)
	filePath := fmt.Sprintf("files/%s.xlsx", fileName)
	invoice.FilePath = filePath

	err = invoice.createXlsxFile(filePath, &invSummary)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to create *.xlsx file: %v", err))
	}

	invoice.Created = time.Now()
	invoice.PeriodFrom = startTime
	invoice.PeriodTo = endTime

	// Store invoice entry
	o := orm.NewOrm()

	/*
		invoiceId, err = o.Insert(&invoice)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to insert invoice into db: %v", err))
		}
		beego.Trace("Created invoice ID:", invoiceId)
	*/

	// Beego time management is very strange...
	// Thinking of converting all datetime fields to string fields in models
	query := fmt.Sprintf("INSERT INTO %s VALUES (?,?,?,?,?,?)",
		invoice.TableName())

	var res sql.Result
	res, err = o.Raw(query,
		nil, invoice.Activations, invoice.FilePath,
		time.Now().Format("2006-01-02 15:04:05"),
		invoice.PeriodFrom.Format("2006-01-02 15:04:05"),
		invoice.PeriodTo.Format("2006-01-02 15:04:05")).Exec()

	if err != nil {
		beego.Error("Failed to insert invoice into db:", err)
		return nil, errors.New(
			fmt.Sprintf("Failed to insert invoice into db: %v", err))
	}

	invoice.Id, err = res.LastInsertId()
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to acquire last inserted id: %v", err))
	}
	return &invoice, nil
}

func GetAllInvoices() (*[]Invoice, error) {

	type CustomInvoice struct {
		Id          int64
		Activations string
		FilePath    string
		Created     string
		PeriodFrom  string
		PeriodTo    string
	}

	customInvoices := []CustomInvoice{}
	inv := Invoice{}
	o := orm.NewOrm()

	query := fmt.Sprintf("SELECT i.* FROM %s i ORDER BY i.id DESC",
		inv.TableName())
	num, err := o.Raw(query).QueryRows(&customInvoices)

	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get all invoices: %v", err))
	}
	beego.Trace("Got num invoices:", num)

	invoices := []Invoice{}

	for invIter := 0; invIter < len(customInvoices); invIter++ {
		inv := Invoice{}
		inv.Id = customInvoices[invIter].Id
		inv.Activations = customInvoices[invIter].Activations
		inv.FilePath = customInvoices[invIter].FilePath
		inv.Created, err = time.ParseInLocation("2006-01-02 15:04:05",
			customInvoices[invIter].Created, time.Now().Location())
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to parse invoice: %v", err))
		}
		inv.PeriodFrom, err = time.ParseInLocation("2006-01-02 15:04:05",
			customInvoices[invIter].PeriodFrom, time.Now().Location())
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to parse invoice: %v", err))
		}
		inv.PeriodTo, err = time.ParseInLocation("2006-01-02 15:04:05",
			customInvoices[invIter].PeriodTo, time.Now().Location())
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to parse invoice: %v", err))
		}

		invoices = append(invoices, inv)
	}

	return &invoices, nil
}

func DeleteInvoice(invoiceId int64) error {
	invoice := Invoice{}
	invoice.Id = invoiceId
	o := orm.NewOrm()
	num, err := o.Delete(&invoice)
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete invoice: %v", err))
	}
	beego.Trace("Deleted num invoices:", num)
	return nil
}

func (this *Invoice) createXlsxFile(filePath string,
	invSummarry *InvoiceSummary) error {

	userSummaries := &(*invSummarry).UserSummaries

	// Create a xlsx file if there
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

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

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "User"
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%s (%s)",
			userSummary.UserFullName, userSummary.Username)
		cell = row.AddCell()
		cell.Value = userSummary.UserEmail

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Debitor Number"
		cell = row.AddCell()
		cell.Value = userSummary.DebitorNumber

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "Activations"

		row = sheet.AddRow()
		_ = row.AddCell()
		cell = row.AddCell()
		cell.Value = "Machine Name"
		cell = row.AddCell()
		cell.Value = "Product ID"
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

		sumTotal := 0.0
		sumTotalDisc := 0.0

		for actIter := 0; actIter < len(userSummary.Activations); actIter++ {
			activation := (*userSummary).Activations[actIter]
			row = sheet.AddRow()
			_ = row.AddCell() // one empty for offset
			cell = row.AddCell()
			cell.Value = activation.MachineName

			cell = row.AddCell()
			cell.Value = activation.MachineProductId

			cell = row.AddCell()
			//beego.Trace(activation.MachineUsage)
			cell.Value = strconv.FormatFloat(activation.MachineUsage,
				'f', 4, 64)

			cell = row.AddCell()
			cell.Value = activation.MachineUsageUnit

			cell = row.AddCell()
			cell.Value = strconv.FormatFloat(activation.MachinePricePerUnit,
				'f', 2, 64)

			// Calculate total
			priceTotal := activation.MachineUsage *
				activation.MachinePricePerUnit
			sumTotal += priceTotal
			cell = row.AddCell()
			cell.Value = strconv.FormatFloat(priceTotal, 'f', 2, 64)

			// Memberships
			membershipStr := ""
			for memIter := 0; memIter < len(activation.Memberships); memIter++ {
				membership := (*activation).Memberships[memIter]
				memStr := fmt.Sprintf("%s (%d%%)",
					membership.ShortName,
					membership.MachinePriceDeduction)
				if membershipStr == "" {
					membershipStr = memStr
				} else {
					membershipStr = fmt.Sprintf("%s, %s",
						membershipStr, memStr)
				}

				// Discount total price
				priceTotal = priceTotal - (priceTotal *
					float64(membership.MachinePriceDeduction) / 100.0)
			}
			sumTotalDisc += priceTotal
			if membershipStr == "" {
				membershipStr = "None"
			}
			cell = row.AddCell()
			cell.Value = membershipStr

			cell = row.AddCell()
			cell.Value = strconv.FormatFloat(priceTotal, 'f', 2, 64)
		} // for userSummary activations

		row = sheet.AddRow()
		_ = row.AddCell()
		_ = row.AddCell()
		_ = row.AddCell()
		_ = row.AddCell()
		_ = row.AddCell()
		cell = row.AddCell()
		cell.Value = "Subtotal €"
		cell = row.AddCell()
		cell.Value = strconv.FormatFloat(sumTotal, 'f', 2, 64)
		cell = row.AddCell()
		cell.Value = "Discounted €"
		cell = row.AddCell()
		cell.Value = strconv.FormatFloat(sumTotalDisc, 'f', 2, 64)

		_ = sheet.AddRow() // extra row after each user summary
	} // for userSummaries

	err = file.Save(filePath)
	if err != nil {
		return err
	}

	this.FilePath = filePath
	return nil
}

func (this *Invoice) getActivations(startTime,
	endTime time.Time) (*[]Activation, error) {

	act := Activation{}
	usr := User{}
	var err error
	o := orm.NewOrm()

	query := fmt.Sprintf("SELECT a.* FROM %s a JOIN %s u ON a.user_id=u.id "+
		"WHERE a.time_start > ? AND a.time_end < ? "+
		"AND a.invoiced = false AND a.active = false "+
		"ORDER BY u.first_name ASC, a.machine_id",
		act.TableName(),
		usr.TableName())

	activations := []Activation{}
	_, err = o.Raw(query,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02")).QueryRows(&activations)
	if err != nil {
		return nil, err
	}

	this.Activations = "["
	for actIter := 0; actIter < len(activations); actIter++ {
		var format string
		if actIter == 0 {
			format = "%s%d"
		} else {
			format = "%s,%d"
		}
		this.Activations = fmt.Sprintf(format,
			this.Activations,
			activations[actIter].Id)
	}
	this.Activations = fmt.Sprintf("%s]", this.Activations)

	return &activations, nil
}

func (this *Invoice) getInvoiceFileName(startTime,
	endTime time.Time) string {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 10)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return fmt.Sprintf("invoice-%s-%s-%s",
		startTime.Format("20060102"),
		endTime.Format("20060102"),
		string(b))
}

func (this *Invoice) getEnhancedActivations(
	activations *[]Activation) (*[]*InvoiceActivation, error) {

	enhActivations := []*InvoiceActivation{}

	// Enhance each activation in the activations slice.
	for actIter := 0; actIter < len(*activations); actIter++ {
		var activation *Activation = &(*activations)[actIter]
		invActivation, err := this.enhanceActivation(activation)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to enhance activation: %v", err))
		}
		enhActivations = append(enhActivations, invActivation)
	}

	return &enhActivations, nil
}

func (this *Invoice) getUserSummaries(
	invoiceActivations *[]*InvoiceActivation) *[]*UserSummary {

	// Create a slice for unique user summaries.
	userSummaries := []*UserSummary{}

	// Sort invoice activations by user.
	for invActIter := 0; invActIter < len(*invoiceActivations); invActIter++ {

		// Search for user Id in the user summaries slice.
		iActivation := (*invoiceActivations)[invActIter]
		uSummaryExists := false
		var summary *UserSummary

		for usrSumIter := 0; usrSumIter < len(userSummaries); usrSumIter++ {
			if iActivation.UserId == userSummaries[usrSumIter].UserId {
				uSummaryExists = true
				summary = userSummaries[usrSumIter]
				break
			}
		}

		// Create new user summary if it does not exist for the user.
		if !uSummaryExists {
			newSummary := UserSummary{}
			newSummary.UserId = iActivation.UserId
			newSummary.UserFullName = iActivation.UserFullName
			newSummary.Username = iActivation.Username
			newSummary.UserEmail = iActivation.UserEmail
			newSummary.DebitorNumber = "Undefined"
			userSummaries = append(userSummaries, &newSummary)
			summary = userSummaries[len(userSummaries)-1]
		}

		// Append the invoice activation to the user summary.
		if summary.UserId == iActivation.UserId {
			summary.Activations = append(summary.Activations, iActivation)
		}
	} // for

	// Return populated user summaries slice.
	return &userSummaries
}

func (this *Invoice) enhanceActivation(activation *Activation) (*InvoiceActivation, error) {

	invActivation := &InvoiceActivation{}
	o := orm.NewOrm()

	// Get activation machine data
	machine := &Machine{}
	query := fmt.Sprintf("SELECT m.name, m.price, m.price_unit FROM %s m "+
		"WHERE id = ?", machine.TableName())

	err := o.Raw(query, activation.MachineId).QueryRow(machine)
	if err != nil {
		//return nil, errors.New(fmt.Sprintf("Failed to get machine: %v", err))
		beego.Info("Failed to get machine, ID: ", activation.MachineId)
	}

	invActivation.MachineId = int64(activation.MachineId)
	invActivation.MachineName = machine.Name
	invActivation.MachineProductId = "Undefined"
	invActivation.MachinePricePerUnit = float64(machine.Price)
	invActivation.MachineUsageUnit = machine.PriceUnit

	// Usage time is stored as seconds and we need to transform that into
	// other format depending on the machine usage unit.
	if invActivation.MachineUsageUnit == "minute" {
		invActivation.MachineUsage = float64(activation.TimeTotal) / 60.0
		if invActivation.MachineUsage < 0.01 {
			invActivation.MachineUsage = 0.01
		}
	} else if invActivation.MachineUsageUnit == "hour" {
		invActivation.MachineUsage = float64(activation.TimeTotal) / 60.0 / 60.0
		if invActivation.MachineUsage < 0.01 {
			invActivation.MachineUsage = 0.01
		}
	}

	// Get activation user data
	user := &User{}
	query = fmt.Sprintf("SELECT first_name, last_name, username, email "+
		"FROM %s WHERE id = ?", user.TableName())
	err = o.Raw(query, activation.UserId).QueryRow(user)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get user: %v", err))
	}

	invActivation.UserId = activation.UserId
	invActivation.UserFullName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	invActivation.Username = user.Username
	invActivation.UserEmail = user.Email
	invActivation.UserDebitorNumber = "Undefined"

	// Get user memberships
	m := &UserMembership{} // Use just for the TableName func
	type CustomUserMembership struct {
		Id           int64 `orm:"auto";"pk"`
		UserId       int64
		MembershipId int64
		StartDate    string
	}
	usrMemberships := &[]CustomUserMembership{}
	query = fmt.Sprintf("SELECT membership_id, start_date FROM %s "+
		"WHERE user_id = ?", m.TableName())
	_, err = o.Raw(query, invActivation.UserId).QueryRows(usrMemberships)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get user membership: %v", err))
	}

	// Check if the membership dates of the user overlap with the activation.
	// If they overlap, add the membership to the invActivation
	for i := 0; i < len(*usrMemberships); i++ {
		//for _, usrMem := range *usrMemberships {

		usrMem := &(*usrMemberships)[i]

		// Parse user membership start time
		var memStartTime time.Time
		memStartTime, err = time.ParseInLocation("2006-01-02 15:04:05",
			usrMem.StartDate,
			time.Now().Location())
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to parse user membership start time: %v",
					err))
		}

		// We also need to get the membership end time that is not yet saved
		// in the database

		// We need to get the actual membership to be able to do that
		mem := &Membership{}
		query = fmt.Sprintf("SELECT title, short_name, duration, unit, "+
			"machine_price_deduction, affected_machines FROM %s "+
			"WHERE id = ?", mem.TableName())
		err = o.Raw(query, usrMem.MembershipId).QueryRow(mem)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to get the actual membership: %v", err))
		}

		// From the membership we need the duration in days to calculate the
		// membership end time.
		var dur time.Duration
		var durStr string
		var durHours int64
		if mem.Unit == "days" {
			durHours = int64(mem.Duration) * 24
			durStr = fmt.Sprintf("%dh", durHours)
		} else {
			return nil, errors.New(
				fmt.Sprintf("Unsupported membership duration unit: %s",
					mem.Unit))
		}

		// Parse the membership duration
		//beego.Trace(durStr)
		dur, err = time.ParseDuration(durStr)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to parse membership duration: %v", err))
		}

		// Calculate end time using user membership start time and
		// membership duration
		var memEndTime time.Time
		memEndTime = memStartTime.Add(dur)

		// Now that we have membership start and end time, let's check
		// if this period of time overlaps with the activation
		if activation.TimeStart.After(memStartTime) &&
			activation.TimeStart.Before(memEndTime) {

			// Yes, this activation is within the range of a membership,
			// add the membership to the activation
			invActivation.Memberships = append(invActivation.Memberships, mem)
		}
	}

	return invActivation, nil
}
