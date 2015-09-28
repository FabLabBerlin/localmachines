package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	orm.RegisterModel(new(Invoice))
}

// Invoice entry that is saved in the database.
// Activations field contains a JSON array with activation IDs.
// XlsFile field contains URL to the generated XLSX file.
type Invoice struct {
	Id            int64  `orm:"auto";"pk"`
	Activations   string `orm:type(text)`
	FilePath      string `orm:size(255)`
	Created       time.Time
	PeriodFrom    time.Time
	PeriodTo      time.Time
	UserSummaries []*UserSummary `orm:"-"`
}

func (this *Invoice) Len() int {
	return len(this.UserSummaries)
}

func (this *Invoice) Less(i, j int) bool {
	a := this.UserSummaries[i]
	b := this.UserSummaries[j]
	aName := a.User.FirstName + " " + a.User.LastName
	bName := b.User.FirstName + " " + b.User.LastName
	return strings.ToLower(aName) < strings.ToLower(bName)
}

func (this *Invoice) Swap(i, j int) {
	this.UserSummaries[i], this.UserSummaries[j] = this.UserSummaries[j], this.UserSummaries[i]
}

func (this *Invoice) TableName() string {
	return "invoices"
}

// This is a user activation row that appears in the XLSX file
type InvoiceActivation struct {
	Activation      Activation
	Machine         *Machine
	MachineUsage    float64
	User            User
	Memberships     []*Membership
	TotalPrice      float64
	DiscountedTotal float64
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

func PriceTotalExclDisc(invAct *InvoiceActivation) float64 {
	return invAct.MachineUsage * float64(invAct.Machine.Price)
}

func PriceTotalDisc(invAct *InvoiceActivation) (float64, error) {
	priceTotal := PriceTotalExclDisc(invAct)
	for _, membership := range invAct.Memberships {

		// We need to know whether the machine is affected by the base membership
		// as well as the individual activation is affected by the user membership
		isAffected, err := membership.IsMachineAffected(invAct.Machine.Id)
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
	if (*this[i]).Activation.TimeStart.Before((*this[j]).Activation.TimeStart) {
		return true
	} else if (*this[j]).Activation.TimeStart.Before((*this[i]).Activation.TimeStart) {
		return false
	} else {
		return (*this[i]).Machine.Name < (*this[j]).Machine.Name
	}
}

func (this InvoiceActivations) Swap(i, j int) {
	*this[i], *this[j] = *this[j], *this[i]
}

func (this InvoiceActivations) SummarizedByMachine() (
	InvoiceActivations, error) {

	byMachine := make(map[string]*InvoiceActivation)
	for _, activation := range this {
		summary, ok := byMachine[activation.Machine.Name]
		if !ok {
			summary = &InvoiceActivation{
				Activation:      Activation{},
				MachineUsage:    0,
				TotalPrice:      0,
				DiscountedTotal: 0,
				Machine:         activation.Machine,
				Memberships:     activation.Memberships,
			}
			byMachine[activation.Machine.Name] = summary
		}
		summary.MachineUsage += activation.MachineUsage
		summary.TotalPrice += activation.TotalPrice
		summary.DiscountedTotal += activation.DiscountedTotal

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
	if (*this[i]).Machine.Name < (*this[j]).Machine.Name {
		return true
	} else if (*this[j]).Machine.Name < (*this[i]).Machine.Name {
		return false
	} else {
		return (*this[i]).Activation.TimeStart.Before((*this[j]).Activation.TimeStart)
	}
}

func (this InvoiceActivationsXlsx) Swap(i, j int) {
	*this[i], *this[j] = *this[j], *this[i]
}

type UserSummary struct {
	User        User
	Activations InvoiceActivations
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Creates invoice entry in the database
func CreateInvoice(startTime, endTime time.Time) (*Invoice, error) {

	var err error

	invoice, err := CalculateInvoiceSummary(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("CalculateInvoiceSummary: %v", err)
	}

	// Create *.xlsx file.
	fileName := invoice.getInvoiceFileName(startTime, endTime)

	// Make sure the files directory exists
	exists, _ := exists("files")
	if !exists {

		// Create the files directory with permission to write
		err = os.Mkdir("files", 0777)
		if err != nil {
			beego.Error("Failed to create files dir:", err)
			return nil, fmt.Errorf("Failed to create files dir: %v", err)
		}
	}

	filePath := fmt.Sprintf("files/%s.xlsx", fileName)
	invoice.FilePath = filePath

	err = createXlsxFile(filePath, &invoice)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to create *.xlsx file: %v", err))
	}

	invoice.Created = time.Now()
	invoice.PeriodFrom = startTime
	invoice.PeriodTo = endTime

	// Store invoice entry
	o := orm.NewOrm()
	invoice.Id, err = o.Insert(&invoice)
	if err != nil {
		beego.Error("Failed to insert invoice into db:", err)
		return nil, fmt.Errorf("Failed to insert invoice into db: %v", err)
	}

	return &invoice, nil
}

// Gets existing invoice from db by invoice ID
func GetInvoice(invoiceId int64) (invoice *Invoice, err error) {

	invoice = &Invoice{}
	invoice.Id = invoiceId

	o := orm.NewOrm()
	err = o.Read(invoice)
	if err != nil {
		beego.Error("Failed to read invoice:", err)
		return nil, fmt.Errorf("Failed to read invoice: %v", err)
	}

	return
}

// Returns Invoice and InvoiceSummary objects, error otherwise
func CalculateInvoiceSummary(startTime, endTime time.Time) (invoice Invoice, err error) {

	// Enhance activations with user and membership data
	var invActivations *[]*InvoiceActivation
	invActivations, err = invoice.getInvoiceActivations(startTime, endTime)
	if err != nil {
		err = fmt.Errorf("Failed to get enhanced activations: %v", err)
		return
	}

	activationIds := make([]string, 0, len(*invActivations))
	for _, act := range *invActivations {
		activationIds = append(activationIds, strconv.FormatInt(act.Activation.Id, 10))
	}
	invoice.Activations = "[" + strings.Join(activationIds, ",") + "]"

	// Create user summaries from invoice activations
	var userSummaries *[]*UserSummary
	userSummaries, err = invoice.getUserSummaries(invActivations)
	if err != nil {
		err = fmt.Errorf("Failed to get user summaries: %v", err)
		return
	}

	for i := 0; i < len(*userSummaries); i++ {
		sort.Stable((*userSummaries)[i].Activations)
		beego.Trace((*userSummaries)[i].Activations)
		for _, activation := range (*userSummaries)[i].Activations {
			activation.TotalPrice = PriceTotalExclDisc(activation)
			activation.DiscountedTotal, err = PriceTotalDisc(activation)
			if err != nil {
				return
			}
		}
	}

	// Create invoice summary
	invoice.PeriodFrom = startTime
	invoice.PeriodTo = endTime
	invoice.UserSummaries = *userSummaries

	return invoice, err
}

// Gets all invoices from the database
func GetAllInvoices() (invoices *[]Invoice, err error) {
	inv := Invoice{}
	var readInvoices []Invoice
	o := orm.NewOrm()
	var num int64
	num, err = o.QueryTable(inv.TableName()).OrderBy("-Id").All(&readInvoices)
	if err != nil {
		beego.Error("Failed to get all invoices:", err)
		return nil, fmt.Errorf("Failed to get all invoices: %v", err)
	}
	beego.Info("Got num invoices:", num)

	return &readInvoices, nil
}

// Deletes an invoice by ID
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

// Gets activations that have happened between start and end dates
func getActivations(startTime,
	endTime time.Time) (activationsArr *[]Activation, err error) {

	act := Activation{}
	usr := User{}
	o := orm.NewOrm()

	query := fmt.Sprintf("SELECT a.* FROM %s a JOIN %s u ON a.user_id=u.id "+
		"WHERE a.time_start > ? AND a.time_end < ? "+
		"AND a.invoiced = false AND a.active = false ",
		act.TableName(),
		usr.TableName())

	activations := []Activation{}
	_, err = o.Raw(query,
		startTime.Format("2006-01-02 15:04:05"),
		endTime.Format("2006-01-02 15:04:05")).QueryRows(&activations)
	if err != nil {
		return nil, err
	}

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

func (this *Invoice) getInvoiceActivations(startTime, endTime time.Time) (*[]*InvoiceActivation, error) {
	// Get all uninvoiced activations in the time range
	var activations *[]Activation
	activations, err := getActivations(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("Failed to get activations: %v", err)
	}

	enhActivations := make([]*InvoiceActivation, 0, len(*activations))

	// Enhance each activation in the activations slice.
	for _, activation := range *activations {
		invActivation, err := this.enhanceActivation(&activation)
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

		for _, userSummary := range userSummaries {
			if iActivation.User.Id == userSummary.User.Id {
				uSummaryExists = true
				summary = userSummary
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

	o := orm.NewOrm()

	// Get activation machine data
	machine, err := GetMachine(activation.MachineId)
	if err != nil {
		beego.Error("Failed to get machine, ID: ", activation.MachineId, ":", err)
		return nil, fmt.Errorf("Failed to get machine: %v", err)
	}

	invActivation := &InvoiceActivation{
		Machine: machine,
	}

	// Usage time is stored as seconds and we need to transform that into
	// other format depending on the machine usage unit.
	switch invActivation.Machine.PriceUnit {
	case "minute":
		invActivation.MachineUsage = float64(activation.TimeTotal) / 60.0
		if invActivation.MachineUsage < 0.01 {
			invActivation.MachineUsage = 0.01
		}
		break
	case "hour":
		invActivation.MachineUsage = float64(activation.TimeTotal) / 60.0 / 60.0
		if invActivation.MachineUsage < 0.01 {
			invActivation.MachineUsage = 0.01
		}
		break
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
	query := fmt.Sprintf("SELECT id, user_id, membership_id, start_date, end_date, auto_extend FROM %s "+
		"WHERE user_id=?", m.TableName())
	_, err = o.Raw(query, invActivation.User.Id).QueryRows(usrMemberships)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get user membership: %v", err))
	}
	// Check if the membership dates of the user overlap with the activation.
	// If they overlap, add the membership to the invActivation
	for i := 0; i < len(*usrMemberships); i++ {
		usrMem := &(*usrMemberships)[i]

		beego.Trace("usrMem.StartTime:", usrMem.StartDate)

		// Get membership
		mem, err := GetMembership(usrMem.MembershipId)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to get the actual membership: %v", err))
		}

		if usrMem.EndDate.IsZero() {
			return nil, fmt.Errorf("end date is zero")
		}

		// Update user membership with correct end date
		userMembershipUpdate := UserMembership{
			Id: usrMem.Id,
		}
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
	invActivation.Activation = *activation
	beego.Info("invActivation.Activation.TimeStart = ", invActivation.Activation.TimeStart)

	return invActivation, nil
}
