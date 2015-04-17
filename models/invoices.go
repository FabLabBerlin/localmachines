package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/tealeg/xlsx"
	"time"
)

type Invoice struct {
	Id          int64  `orm:"auto";"pk"`
	Activations string `orm:type(text)`
	XlsFile     string `orm:size(255)`
}

func (this *Invoice) TableName() string {
	return "invoices"
}

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

func init() {
	orm.RegisterModel(new(Invoice))
}

func CreateInvoice(startTime, endTime time.Time) (*Invoice, error) {

	beego.Info("Creating invoice...")

	var invoice *Invoice

	var err error

	// Create invoice object that will be returned
	invoice = &Invoice{}

	// Get all uninvoiced activations in the time range
	var activations *[]Activation
	activations, err = invoice.getActivations(startTime, endTime)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get activations: %v", err))
	}

	// Enhance activations with user and membership data
	var invActivations *[]InvoiceActivation
	invActivations, err = invoice.getEnhancedActivations(activations)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get enhanced activations: %v", err))
	}

	// Create user summaries from invoice activations
	var userSummaries *[]UserSummary
	userSummaries = invoice.getUserSummaries(invActivations)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get user summaries: %v", err))
	}

	for i := 0; i < len(*userSummaries); i++ {
		beego.Trace((*userSummaries)[i].Activations)
	}

	// Create a xlsx file if there
	var file *xlsx.File
	//var sheet *xlsx.Sheet
	//var row *xlsx.Row
	//var cell *xlsx.Cell
	file = xlsx.NewFile()
	//sheet = file.AddSheet("Sheet1")

	// Fill the xlsx sheet

	filename := "files/MyXLSXFile.xlsx"
	err = file.Save(filename)
	if err != nil {
		fmt.Printf(err.Error())
	}

	inv := Invoice{}
	inv.XlsFile = filename

	return &inv, nil
}

func (this *Invoice) getActivations(startTime,
	endTime time.Time) (*[]Activation, error) {

	act := Activation{}
	usr := User{}
	var err error
	o := orm.NewOrm()

	query := fmt.Sprintf("SELECT * FROM %s a JOIN %s u ON a.user_id=u.id "+
		"WHERE a.time_start > ? AND a.time_end < ? "+
		"AND a.invoiced = false AND a.active = false "+
		"ORDER BY u.first_name ASC",
		act.TableName(),
		usr.TableName())

	activations := []Activation{}
	_, err = o.Raw(query,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02")).QueryRows(&activations)
	if err != nil {
		return nil, err
	}

	return &activations, nil
}

func (this *Invoice) getEnhancedActivations(
	activations *[]Activation) (*[]InvoiceActivation, error) {

	enhancedActivations := []InvoiceActivation{}
	for _, a := range *activations {
		var ia *InvoiceActivation
		var err error
		ia, err = this.enhanceActivation(&a)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to enhance activation: %v", err))
		}
		enhancedActivations = append(enhancedActivations, *ia)
	}
	return &enhancedActivations, nil
}

func (this *Invoice) getUserSummaries(
	invoiceActivations *[]InvoiceActivation) *[]UserSummary {

	// Create a slice for unique user summaries
	userSummaries := []UserSummary{}

	// Sort invoice activations by user
	//for _, iActivation := range *invoiceActivations {
	for j := 0; j < len(*invoiceActivations); j++ {
		// Search for user Id in the user summaries slice

		iActivation := &(*invoiceActivations)[j]

		uSummaryExists := false
		var summary *UserSummary

		for i := 0; i < len(userSummaries); i++ {
			if iActivation.UserId == userSummaries[i].UserId {
				uSummaryExists = true
				beego.Trace("Summary exists")
				summary = &userSummaries[i]
				break
			}
		}

		if !uSummaryExists {
			beego.Trace("Creating new summary")
			summary = &UserSummary{}
			summary.UserId = iActivation.UserId
			summary.UserFullName = iActivation.UserFullName
			summary.Username = iActivation.Username
			summary.UserEmail = iActivation.UserEmail
			summary.DebitorNumber = "Undefined"

			// For some reason the first activation is not appended
			summary.Activations = []*InvoiceActivation{iActivation}

			userSummaries = append(userSummaries, *summary)
			beego.Trace("Created new user summary for user", summary.UserFullName)
		}

		// Append the invoice activation to the user summary
		if summary.UserId == iActivation.UserId {
			beego.Trace("Appending activation")
			summary.Activations = append(summary.Activations, iActivation)
			beego.Trace(summary.Activations)
		}
	} // for

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
		query = fmt.Sprintf("SELECT title, duration, unit, "+
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
