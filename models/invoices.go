package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"sort"
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
	User                User
	Memberships         []*Membership
	TotalPrice          float64
	DiscountedTotal     float64
	TimeStart           time.Time
	TimeEnd             time.Time
}

func (this *InvoiceActivation) MembershipStr() string {
	membershipStr := ""
	for _, membership := range this.Memberships {
		memStr := fmt.Sprintf("%s (%d%%)",
			membership.ShortName,
			membership.MachinePriceDeduction)
		if membershipStr == "" {
			membershipStr = memStr
		} else {
			membershipStr = fmt.Sprintf("%s, %s",
				membershipStr, memStr)
		}
	}
	if membershipStr == "" {
		membershipStr = "None"
	}
	return membershipStr
}

func (this *InvoiceActivation) PriceTotalExclDisc() float64 {
	return this.MachineUsage * this.MachinePricePerUnit
}

func (this *InvoiceActivation) PriceTotalDisc() (float64, error) {
	priceTotal := this.PriceTotalExclDisc()
	for _, membership := range this.Memberships {

		// We need to know whether the machine is affected by the base membership
		// as well as the individual activation is affected by the user membership
		isAffected, err := membership.IsMachineAffected(this.MachineId)
		if err != nil {
			beego.Error(
				"Failed to check whether machine is affected by membership:", err)
			return 0, fmt.Errorf(
				"Failed to check whether machine is affected by membership")
		}

		machinePriceDeduction := 0.0
		if isAffected {
			machinePriceDeduction = float64(membership.MachinePriceDeduction)
		}
		// Discount total price
		priceTotal = priceTotal - (priceTotal * machinePriceDeduction / 100.0)
	}
	return priceTotal, nil
}

type InvoiceActivations []*InvoiceActivation

func (this InvoiceActivations) Len() int {
	return len(this)
}

func (this InvoiceActivations) Less(i, j int) bool {
	if (*this[i]).TimeStart.Before((*this[j]).TimeStart) {
		return true
	} else if (*this[j]).TimeStart.Before((*this[i]).TimeStart) {
		return false
	} else {
		return (*this[i]).MachineName < (*this[j]).MachineName
	}
}

func (this InvoiceActivations) Swap(i, j int) {
	*this[i], *this[j] = *this[j], *this[i]
}

func (this InvoiceActivations) SummarizedByMachine() (InvoiceActivations, error) {
	byMachine := make(map[string]*InvoiceActivation)
	for _, activation := range this {
		summary, ok := byMachine[activation.MachineName]
		if !ok {
			a := *activation
			summary = &a
			summary.TimeStart = time.Unix(0, 0)
			summary.MachineUsage = 0
			summary.TotalPrice = 0
			summary.DiscountedTotal = 0
			byMachine[activation.MachineName] = summary
		}
		summary.MachineUsage += activation.MachineUsage
		summary.TotalPrice += activation.PriceTotalExclDisc()
		if priceTotalDisc, err := activation.PriceTotalDisc(); err == nil {
			summary.DiscountedTotal += priceTotalDisc
		} else {
			return nil, err
		}

	}

	sumActivations := make(InvoiceActivations, 0, len(byMachine))
	for _, summary := range byMachine {
		sumActivations = append(sumActivations, summary)
	}
	sort.Stable(sumActivations)

	return sumActivations, nil
}

type InvoiceActivationsXlsx []*InvoiceActivation

func (this InvoiceActivationsXlsx) Len() int {
	return len(this)
}

func (this InvoiceActivationsXlsx) Less(i, j int) bool {
	if (*this[i]).MachineName < (*this[j]).MachineName {
		return true
	} else if (*this[j]).MachineName < (*this[i]).MachineName {
		return false
	} else {
		return (*this[i]).TimeStart.Before((*this[j]).TimeStart)
	}
}

func (this InvoiceActivationsXlsx) Swap(i, j int) {
	*this[i], *this[j] = *this[j], *this[i]
}

type UserSummary struct {
	User        User
	Activations InvoiceActivations
}

type InvoiceSummary struct {
	InvoiceName     string
	PeriodStartTime time.Time
	PeriodEndTime   time.Time
	UserSummaries   []*UserSummary
}

func CreateInvoice(startTime, endTime time.Time) (*Invoice, error) {

	var err error

	invoice, invSummary, err := CalculateInvoiceSummary(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("CalculateInvoiceSummary: %v", err)
	}

	// Create *.xlsx file.
	fileName := invoice.getInvoiceFileName(startTime, endTime)
	filePath := fmt.Sprintf("files/%s.xlsx", fileName)
	invoice.FilePath = filePath

	err = createXlsxFile(filePath, &invoice, &invSummary)
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

func CalculateInvoiceSummary(
	startTime, endTime time.Time) (
	invoice Invoice, invSummary InvoiceSummary, err error) {

	// Get all uninvoiced activations in the time range
	var activations *[]Activation
	activations, err = invoice.getActivations(startTime, endTime)
	if err != nil {
		err = fmt.Errorf("Failed to get activations: %v", err)
		return
	}

	beego.Trace("Activations str:", invoice.Activations)

	// Enhance activations with user and membership data
	var enhancedActivations *[]*InvoiceActivation
	enhancedActivations, err = invoice.getEnhancedActivations(activations)
	if err != nil {
		err = fmt.Errorf("Failed to get enhanced activations: %v", err)
		return
	}

	// Create user summaries from invoice activations
	var userSummaries *[]*UserSummary
	userSummaries, err = invoice.getUserSummaries(enhancedActivations)
	if err != nil {
		err = fmt.Errorf("Failed to get user summaries: %v", err)
		return
	}

	for i := 0; i < len(*userSummaries); i++ {
		sort.Stable((*userSummaries)[i].Activations)
		beego.Trace((*userSummaries)[i].Activations)
		for _, activation := range (*userSummaries)[i].Activations {
			activation.TotalPrice = activation.PriceTotalExclDisc()
			activation.DiscountedTotal, err = activation.PriceTotalDisc()
			if err != nil {
				return
			}
		}
	}

	// Create invoice summary
	invSummary.InvoiceName = "Invoice"
	invSummary.PeriodStartTime = startTime
	invSummary.PeriodEndTime = endTime
	invSummary.UserSummaries = *userSummaries

	return invoice, invSummary, err
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

		// For each activation get the correct start and end time
		type ActivationTimes struct {
			TimeStart string
			TimeEnd   string
		}
		activationTimes := ActivationTimes{}
		query = fmt.Sprintf(
			"SELECT time_start, time_end FROM %s WHERE id=?", act.TableName())
		err = o.Raw(query, activations[actIter].Id).QueryRow(&activationTimes)
		if err != nil {
			beego.Error("Could not get activation start and end time as string:", err)
			return nil, fmt.Errorf(
				"Could not get activation start and end time as string")
		}
		beego.Trace("Activation start time:", activationTimes.TimeStart)
		beego.Trace("Activation end time:", activationTimes.TimeEnd)

		// Parse the time strings correctly.
		// For now the time is decoded taking into account that the time
		// has been saved as if for the location of the server in Berlin.
		// In the future the time should be saved as UTC time. This should be
		// done with a clever migration.
		parseLayout := "2006-01-02 15:04:05"
		activations[actIter].TimeStart, _ = time.ParseInLocation(parseLayout,
			activationTimes.TimeStart, time.Now().Location())
		activations[actIter].TimeEnd, _ = time.ParseInLocation(parseLayout,
			activationTimes.TimeEnd, time.Now().Location())

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
	invoiceActivations *[]*InvoiceActivation) (*[]*UserSummary, error) {

	// Create a slice for unique user summaries.
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}
	userSummaries := make([]*UserSummary, 0, len(users))
	for _, user := range users {
		newSummary := UserSummary{}
		newSummary.User = *user
		userSummaries = append(userSummaries, &newSummary)
	}

	// Sort invoice activations by user.
	for invActIter := 0; invActIter < len(*invoiceActivations); invActIter++ {

		// Search for user Id in the user summaries slice.
		iActivation := (*invoiceActivations)[invActIter]
		uSummaryExists := false
		var summary *UserSummary

		for usrSumIter := 0; usrSumIter < len(userSummaries); usrSumIter++ {
			if iActivation.User.Id == userSummaries[usrSumIter].User.Id {
				uSummaryExists = true
				summary = userSummaries[usrSumIter]
				break
			}
		}

		// Create new user summary if it does not exist for the user.
		if !uSummaryExists {
			beego.Warn("Creating user summary for activation that has no matching user")
			newSummary := UserSummary{}
			newSummary.User = iActivation.User
			userSummaries = append(userSummaries, &newSummary)
			summary = userSummaries[len(userSummaries)-1]
		}

		// Append the invoice activation to the user summary.
		if summary.User.Id == iActivation.User.Id {
			summary.Activations = append(summary.Activations, iActivation)
		}
	} // for

	// Return populated user summaries slice.
	return &userSummaries, nil
}

func (this *Invoice) enhanceActivation(activation *Activation) (
	*InvoiceActivation, error) {

	invActivation := &InvoiceActivation{}
	o := orm.NewOrm()

	// Get activation machine data
	machine := &Machine{}
	query := fmt.Sprintf("SELECT m.name, m.price, m.price_unit FROM %s m "+
		"WHERE id = ?", machine.TableName())

	err := o.Raw(query, activation.MachineId).QueryRow(machine)
	if err != nil {
		beego.Error("Failed to get machine, ID: ", activation.MachineId, ":", err)
		return nil, fmt.Errorf("Failed to get machine: %v", err)
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
	user, err := GetUser(activation.UserId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get user: %v", err))
	}

	invActivation.User = *user

	// Get user memberships
	m := &UserMembership{} // Use just for the TableName func
	usrMemberships := &[]UserMembership{}
	query = fmt.Sprintf("SELECT id, membership_id, start_date, end_date FROM %s "+
		"WHERE user_id=?", m.TableName())
	_, err = o.Raw(query, invActivation.User.Id).QueryRows(usrMemberships)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get user membership: %v", err))
	}
	// Check if the membership dates of the user overlap with the activation.
	// If they overlap, add the membership to the invActivation
	for i := 0; i < len(*usrMemberships); i++ {
		//for _, usrMem := range *usrMemberships {

		usrMem := &(*usrMemberships)[i]

		beego.Trace("usrMem.StartTime:", usrMem.StartDate)

		// Get membership
		mem := &Membership{}
		query = fmt.Sprintf("SELECT title, short_name, duration_months, "+
			"machine_price_deduction, affected_machines FROM %s "+
			"WHERE id=?", mem.TableName())
		err = o.Raw(query, usrMem.MembershipId).QueryRow(mem)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to get the actual membership: %v", err))
		}

		if usrMem.EndDate.IsZero() {
			return nil, fmt.Errorf("end date is zero")
		}

		// Update user membership with correct end date
		userMembershipUpdate := UserMembership{}
		userMembershipUpdate.Id = usrMem.Id
		beego.Trace("userMembershipUpdate.Id:", userMembershipUpdate.Id)
		userMembershipUpdate.EndDate = usrMem.EndDate
		_, err = o.Update(&userMembershipUpdate, "EndDate")
		if err != nil {
			beego.Error("Failed to update user membership end date:", err)
			return nil, fmt.Errorf("Failed to update user membership end date")
		}

		// Now that we have membership start and end time, let's check
		// if this period of time overlaps with the activation
		if activation.TimeStart.After(usrMem.StartDate) &&
			activation.TimeStart.Before(usrMem.EndDate) {

			// Yes, this activation is within the range of a membership,
			// add the membership to the activation
			beego.Trace("Activation affected by membership")
			invActivation.Memberships = append(invActivation.Memberships, mem)
		}
	}

	invActivation.TimeStart = activation.TimeStart
	invActivation.TimeEnd = activation.TimeEnd

	return invActivation, nil
}
