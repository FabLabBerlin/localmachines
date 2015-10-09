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

// This is a purchase row that appears in the XLSX file
type Purchase struct {
	Activation      Activation
	Machine         *Machine
	MachineUsage    float64
	User            User
	Memberships     []*Membership
	TotalPrice      float64
	DiscountedTotal float64
}

func (this *Purchase) MembershipStr() string {
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

func PriceTotalExclDisc(p *Purchase) float64 {
	return p.MachineUsage * float64(p.Machine.Price)
}

func PriceTotalDisc(p *Purchase) (float64, error) {
	priceTotal := PriceTotalExclDisc(p)
	for _, membership := range p.Memberships {

		// We need to know whether the machine is affected by the base membership
		// as well as the individual activation is affected by the user membership
		isAffected, err := membership.IsMachineAffected(p.Machine.Id)
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

type Purchases []*Purchase

func (this Purchases) Len() int {
	return len(this)
}

func (this Purchases) Less(i, j int) bool {
	if (*this[i]).Activation.TimeStart.Before((*this[j]).Activation.TimeStart) {
		return true
	} else if (*this[j]).Activation.TimeStart.Before((*this[i]).Activation.TimeStart) {
		return false
	} else {
		return (*this[i]).Machine.Name < (*this[j]).Machine.Name
	}
}

func (this Purchases) Swap(i, j int) {
	*this[i], *this[j] = *this[j], *this[i]
}

func (this Purchases) SummarizedByMachine() (
	Purchases, error) {

	byMachine := make(map[string]*Purchase)
	for _, activation := range this {
		summary, ok := byMachine[activation.Machine.Name]
		if !ok {
			summary = &Purchase{
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

	sumPurchases := make(Purchases, 0, len(byMachine))
	for _, summary := range byMachine {
		sumPurchases = append(sumPurchases, summary)
	}
	sort.Stable(sumPurchases)

	return sumPurchases, nil
}

type UserSummary struct {
	User      User
	Purchases Purchases
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
	var purchases []*Purchase
	purchases, err = invoice.getPurchases(startTime, endTime)
	if err != nil {
		err = fmt.Errorf("Failed to get enhanced activations: %v", err)
		return
	}

	activationIds := make([]string, 0, len(purchases))
	for _, act := range purchases {
		activationIds = append(activationIds, strconv.FormatInt(act.Activation.Id, 10))
	}
	invoice.Activations = "[" + strings.Join(activationIds, ",") + "]"

	// Create user summaries from invoice activations
	var userSummaries *[]*UserSummary
	userSummaries, err = invoice.getUserSummaries(purchases)
	if err != nil {
		err = fmt.Errorf("Failed to get user summaries: %v", err)
		return
	}

	for i := 0; i < len(*userSummaries); i++ {
		sort.Stable((*userSummaries)[i].Purchases)
		beego.Trace((*userSummaries)[i].Purchases)
		for _, activation := range (*userSummaries)[i].Purchases {
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

func (this *Invoice) getPurchases(startTime, endTime time.Time) ([]*Purchase, error) {
	// Get all uninvoiced activations in the time range
	var activations *[]Activation
	activations, err := getActivations(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("Failed to get activations: %v", err)
	}

	enhActivations := make([]*Purchase, 0, len(*activations))

	machines, err := GetAllMachines()
	if err != nil {
		return nil, fmt.Errorf("Failed to get machines: %v", err)
	}
	machinesById := make(map[int64]*Machine)
	for _, machine := range machines {
		machinesById[machine.Id] = machine
	}

	users, err := GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("Failed to get users: %v", err)
	}
	usersById := make(map[int64]User)
	for _, user := range users {
		usersById[user.Id] = *user
	}

	userMemberships, err := GetAllUserMemberships()
	if err != nil {
		return nil, fmt.Errorf("Failed to get user memberships: %v", err)
	}
	userMembershipsById := make(map[int64][]*UserMembership)
	for _, userMembership := range userMemberships {
		uid := userMembership.UserId
		if _, ok := userMembershipsById[uid]; !ok {
			userMembershipsById[uid] = []*UserMembership{
				userMembership,
			}
		} else {
			userMembershipsById[uid] = append(userMembershipsById[uid], userMembership)
		}
	}

	memberships, err := GetAllMemberships()
	if err != nil {
		return nil, fmt.Errorf("Failed to get memberships: %v", err)
	}
	membershipsById := make(map[int64]*Membership)
	for _, membership := range memberships {
		membershipsById[membership.Id] = membership
	}

	// Enhance each activation in the activations slice.
	for _, activation := range *activations {
		invActivation, err := this.enhanceActivation(&activation, machinesById, usersById, userMembershipsById, membershipsById)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to enhance activation: %v", err))
		}
		enhActivations = append(enhActivations, invActivation)
	}

	return enhActivations, nil
}

func (this *Invoice) getUserSummaries(
	purchases []*Purchase) (*[]*UserSummary, error) {

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

	// Sort purchases by user.
	for _, purchase := range purchases {

		uSummaryExists := false
		var summary *UserSummary

		for _, userSummary := range userSummaries {
			if purchase.User.Id == userSummary.User.Id {
				uSummaryExists = true
				summary = userSummary
				break
			}
		}

		// Create new user summary if it does not exist for the user.
		if !uSummaryExists {
			beego.Warn("Creating user summary for activation that has no matching user")
			newSummary := UserSummary{}
			newSummary.User = purchase.User
			userSummaries = append(userSummaries, &newSummary)
			summary = userSummaries[len(userSummaries)-1]
		}

		// Append the invoice activation to the user summary.
		if summary.User.Id == purchase.User.Id {
			summary.Purchases = append(summary.Purchases, purchase)
		}
	}

	// Return populated user summaries slice.
	return &userSummaries, nil
}

func (this *Invoice) enhanceActivation(activation *Activation, machinesById map[int64]*Machine, usersById map[int64]User, userMembershipsByUserId map[int64][]*UserMembership, membershipsById map[int64]*Membership) (
	*Purchase, error) {

	machine, ok := machinesById[activation.MachineId]
	if !ok {
		return nil, fmt.Errorf("No machine has the ID %v", activation.MachineId)
	}

	purchase := &Purchase{
		Machine: machine,
	}

	// Usage time is stored as seconds and we need to transform that into
	// other format depending on the machine usage unit.
	switch purchase.Machine.PriceUnit {
	case "minute":
		purchase.MachineUsage = float64(activation.TimeTotal) / 60.0
		if purchase.MachineUsage < 0.01 {
			purchase.MachineUsage = 0.01
		}
		break
	case "hour":
		purchase.MachineUsage = float64(activation.TimeTotal) / 60.0 / 60.0
		if purchase.MachineUsage < 0.01 {
			purchase.MachineUsage = 0.01
		}
		break
	}

	if purchase.User, ok = usersById[activation.UserId]; !ok {
		return nil, fmt.Errorf("No user has the ID %v", activation.MachineId)
	}

	usrMemberships, ok := userMembershipsByUserId[activation.UserId]
	if !ok {
		usrMemberships = []*UserMembership{}
	}

	// Check if the membership dates of the user overlap with the activation.
	// If they overlap, add the membership to the invActivation
	for _, usrMem := range usrMemberships {

		// Get membership
		mem, ok := membershipsById[usrMem.MembershipId]
		if !ok {
			return nil, fmt.Errorf("Unknown membership id: %v", usrMem.MembershipId)
		}

		if usrMem.EndDate.IsZero() {
			return nil, fmt.Errorf("end date is zero")
		}

		// Now that we have membership start and end time, let's check
		// if this period of time overlaps with the activation
		if activation.TimeStart.After(usrMem.StartDate) &&
			activation.TimeStart.Before(usrMem.EndDate) {

			purchase.Memberships = append(purchase.Memberships, mem)
		}
	}
	purchase.Activation = *activation

	return purchase, nil
}
