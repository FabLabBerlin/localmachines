package invoices

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/purchases"
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

type UserSummary struct {
	User      models.User
	Purchases purchases.Purchases
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
func Create(startTime, endTime time.Time) (*Invoice, error) {

	var err error

	invoice, err := CalculateSummary(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("CalculateInvoiceSummary: %v", err)
	}

	// Create *.xlsx file.
	fileName := invoice.getFileName(startTime, endTime)

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
func Get(invoiceId int64) (invoice *Invoice, err error) {

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
func CalculateSummary(startTime, endTime time.Time) (invoice Invoice, err error) {
	// Enhance activations with user and membership data
	ps, err := invoice.getPurchases(startTime, endTime)
	if err != nil {
		err = fmt.Errorf("Failed to get enhanced activations: %v", err)
		return
	}

	activationIds := make([]string, 0, len(ps))
	for _, p := range ps {
		if p.Activation != nil {
			activationIds = append(activationIds, strconv.FormatInt(p.Id, 10))
		}
	}
	invoice.Activations = "[" + strings.Join(activationIds, ",") + "]"

	// Create user summaries from invoice activations
	var userSummaries *[]*UserSummary
	userSummaries, err = invoice.getUserSummaries(purchases.Purchases{
		Data: ps,
	})
	if err != nil {
		err = fmt.Errorf("Failed to get user summaries: %v", err)
		return
	}

	for i := 0; i < len(*userSummaries); i++ {
		sort.Stable((*userSummaries)[i].Purchases)
		for _, activation := range (*userSummaries)[i].Purchases.Data {
			activation.TotalPrice = purchases.PriceTotalExclDisc(activation)
			activation.DiscountedTotal, err = purchases.PriceTotalDisc(activation)
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
func GetAll() (invoices *[]Invoice, err error) {
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
func Delete(invoiceId int64) error {
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

// Gets purchases that have happened between start and end dates
func getPurchases(startTime,
	endTime time.Time) (ps []*purchases.Purchase, err error) {

	p := purchases.Purchase{}
	usr := models.User{}
	o := orm.NewOrm()

	query := fmt.Sprintf("SELECT p.* FROM %s p JOIN %s u ON p.user_id=u.id "+
		"WHERE p.time_start > ? AND p.time_end < ? "+
		"AND (p.running IS NULL OR p.running = 0)",
		p.TableName(),
		usr.TableName())

	_, err = o.Raw(query,
		startTime.Format("2006-01-02 15:04:05"),
		endTime.Format("2006-01-02 15:04:05")).QueryRows(&ps)

	return
}

func (this *Invoice) getFileName(startTime, endTime time.Time) string {

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

func (this *Invoice) getPurchases(startTime, endTime time.Time) (ps []*purchases.Purchase, err error) {
	machines, err := models.GetAllMachines(true)
	if err != nil {
		return nil, fmt.Errorf("Failed to get machines: %v", err)
	}
	machinesById := make(map[int64]*models.Machine)
	for _, machine := range machines {
		machinesById[machine.Id] = machine
	}

	users, err := models.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("Failed to get users: %v", err)
	}
	usersById := make(map[int64]models.User)
	for _, user := range users {
		usersById[user.Id] = *user
	}

	userMemberships, err := models.GetAllUserMemberships()
	if err != nil {
		return nil, fmt.Errorf("Failed to get user memberships: %v", err)
	}
	userMembershipsById := make(map[int64][]*models.UserMembership)
	for _, userMembership := range userMemberships {
		uid := userMembership.UserId
		if _, ok := userMembershipsById[uid]; !ok {
			userMembershipsById[uid] = []*models.UserMembership{
				userMembership,
			}
		} else {
			userMembershipsById[uid] = append(userMembershipsById[uid], userMembership)
		}
	}

	memberships, err := models.GetAllMemberships()
	if err != nil {
		return nil, fmt.Errorf("Failed to get memberships: %v", err)
	}
	membershipsById := make(map[int64]*models.Membership)
	for _, membership := range memberships {
		membershipsById[membership.Id] = membership
	}

	// Get all uninvoiced purchases in the time range
	ps, err = getPurchases(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("Failed to get purchases: %v", err)
	}

	// Enhance purchases
	for _, p := range ps {
		err := this.enhancePurchase(p, machinesById,
			usersById, userMembershipsById, membershipsById)
		if err != nil {
			return nil, fmt.Errorf("Failed to enhance purchase: %v", err)
		}
	}

	return
}

func (this *Invoice) getUserSummaries(
	ps purchases.Purchases) (*[]*UserSummary, error) {

	// Create a slice for unique user summaries.
	users, err := models.GetAllUsers()
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
	for _, p := range ps.Data {

		uSummaryExists := false
		var summary *UserSummary

		for _, userSummary := range userSummaries {
			if p.User.Id == userSummary.User.Id {
				uSummaryExists = true
				summary = userSummary
				break
			}
		}

		// Create new user summary if it does not exist for the user.
		if !uSummaryExists {
			beego.Warn("Creating user summary for activation that has no matching user")
			newSummary := UserSummary{}
			newSummary.User = p.User
			userSummaries = append(userSummaries, &newSummary)
			summary = userSummaries[len(userSummaries)-1]
		}

		// Append the invoice activation to the user summary.
		if summary.User.Id == p.User.Id {
			summary.Purchases.Data = append(summary.Purchases.Data, p)
		}
	}

	// Return populated user summaries slice.
	return &userSummaries, nil
}

func (this *Invoice) enhancePurchase(purchase *purchases.Purchase,
	machinesById map[int64]*models.Machine, usersById map[int64]models.User,
	userMembershipsByUserId map[int64][]*models.UserMembership,
	membershipsById map[int64]*models.Membership) error {

	var ok bool
	purchase.Machine, ok = machinesById[purchase.MachineId]
	if !ok && (purchase.Type == purchases.TYPE_ACTIVATION || purchase.Type == purchases.TYPE_RESERVATION) {
		return fmt.Errorf("No machine has the ID %v", purchase.MachineId)
	}

	if purchase.User, ok = usersById[purchase.UserId]; !ok {
		return fmt.Errorf("No user has the ID %v", purchase.MachineId)
	}

	usrMemberships, ok := userMembershipsByUserId[purchase.UserId]
	if !ok {
		usrMemberships = []*models.UserMembership{}
	}

	// Check if the membership dates of the user overlap with the activation.
	// If they overlap, add the membership to the invActivation
	for _, usrMem := range usrMemberships {

		// Get membership
		mem, ok := membershipsById[usrMem.MembershipId]
		if !ok {
			return fmt.Errorf("Unknown membership id: %v", usrMem.MembershipId)
		}

		if usrMem.EndDate.IsZero() {
			return fmt.Errorf("end date is zero")
		}

		// Now that we have membership start and end time, let's check
		// if this period of time overlaps with the activation
		if purchase.TimeStart.After(usrMem.StartDate) &&
			purchase.TimeStart.Before(usrMem.EndDate) {

			purchase.Memberships = append(purchase.Memberships, mem)
		}
	}

	return nil
}