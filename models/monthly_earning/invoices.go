package monthly_earning

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/users"
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
	orm.RegisterModel(new(MonthlyEarning))
}

// MonthlyEarning that is saved in the database.
// Activations field contains a JSON array with activation IDs.
// XlsFile field contains URL to the generated XLSX file.
type MonthlyEarning struct {
	Id          int64 `orm:"auto";"pk"`
	MonthFrom   int
	YearFrom    int
	MonthTo     int
	YearTo      int
	Activations string `orm:type(text)`
	FilePath    string `orm:size(255)`
	Created     time.Time
	Invoices    []*Invoice `orm:"-"`
}

func (this *MonthlyEarning) Interval() lib.Interval {
	return lib.Interval{
		MonthFrom: this.MonthFrom,
		YearFrom:  this.YearFrom,
		MonthTo:   this.MonthTo,
		YearTo:    this.YearTo,
	}
}

func (this *MonthlyEarning) Len() int {
	return len(this.Invoices)
}

func (this *MonthlyEarning) Less(i, j int) bool {
	a := this.Invoices[i]
	b := this.Invoices[j]
	aName := a.User.FirstName + " " + a.User.LastName
	bName := b.User.FirstName + " " + b.User.LastName
	return strings.ToLower(aName) < strings.ToLower(bName)
}

func (this *MonthlyEarning) PeriodFrom() time.Time {
	return this.Interval().TimeFrom()
}

func (this *MonthlyEarning) PeriodTo() time.Time {
	return this.Interval().TimeTo()
}

func (this *MonthlyEarning) Swap(i, j int) {
	this.Invoices[i], this.Invoices[j] = this.Invoices[j], this.Invoices[i]
}

func (this *MonthlyEarning) TableName() string {
	return "monthly_earnings"
}

type Invoice struct {
	User      users.User
	Purchases purchases.Purchases
}

func (inv *Invoice) byProductNameAndPricePerUnit() map[string]map[float64][]*purchases.Purchase {
	byProductNameAndPricePerUnit := make(map[string]map[float64][]*purchases.Purchase)
	for _, p := range inv.Purchases.Data {
		if _, ok := byProductNameAndPricePerUnit[p.ProductName()]; !ok {
			byProductNameAndPricePerUnit[p.ProductName()] = make(map[float64][]*purchases.Purchase)
		}
		if _, ok := byProductNameAndPricePerUnit[p.ProductName()][p.PricePerUnit]; !ok {
			byProductNameAndPricePerUnit[p.ProductName()][p.PricePerUnit] = make([]*purchases.Purchase, 0, 20)
		}
		byProductNameAndPricePerUnit[p.ProductName()][p.PricePerUnit] = append(byProductNameAndPricePerUnit[p.ProductName()][p.PricePerUnit], p)
	}
	return byProductNameAndPricePerUnit
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

// Returns MonthlyEarning with populated Invoices
func New(locationId int64, interval lib.Interval) (me MonthlyEarning, err error) {
	// Enhance activations with user and membership data
	ps, err := me.getPurchases(locationId, interval)
	if err != nil {
		err = fmt.Errorf("Failed to get enhanced purchaes: %v", err)
		return
	}

	activationIds := make([]string, 0, len(ps))
	for _, p := range ps {
		if p.Activation != nil {
			activationIds = append(activationIds, strconv.FormatInt(p.Id, 10))
		}
	}
	me.Activations = "[" + strings.Join(activationIds, ",") + "]"

	// Create invoices from purchases
	invs, err := me.GetInvoices(purchases.Purchases{
		Data: ps,
	})
	if err != nil {
		err = fmt.Errorf("Failed to get user summaries: %v", err)
		return
	}

	for i := 0; i < len(invs); i++ {
		sort.Stable(invs[i].Purchases)
		for _, purchase := range invs[i].Purchases.Data {
			purchase.TotalPrice = purchases.PriceTotalExclDisc(purchase)
			purchase.DiscountedTotal, err = purchases.PriceTotalDisc(purchase)
			if err != nil {
				return
			}
		}
	}

	me.MonthFrom = interval.MonthFrom
	me.YearFrom = interval.YearFrom
	me.MonthTo = interval.MonthTo
	me.YearTo = interval.YearTo
	me.Invoices = invs

	return me, err
}

// Creates monthly earning entry in the database
func Create(locationId int64, interval lib.Interval) (*MonthlyEarning, error) {

	var err error

	me, err := New(locationId, interval)
	if err != nil {
		return nil, fmt.Errorf("New: %v", err)
	}

	// Create *.xlsx file.
	fileName := me.getFileName(interval)

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
	me.FilePath = filePath

	err = createXlsxFile(filePath, &me)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to create *.xlsx file: %v", err))
	}

	me.Created = time.Now()
	me.MonthFrom = interval.MonthFrom
	me.YearFrom = interval.YearFrom
	me.MonthTo = interval.MonthTo
	me.MonthTo = interval.MonthTo

	// Store monthly earning entry
	o := orm.NewOrm()
	me.Id, err = o.Insert(&me)
	if err != nil {
		beego.Error("Failed to insert monthly earning into db:", err)
		return nil, fmt.Errorf("Failed to insert monthly earning into db: %v", err)
	}

	return &me, nil
}

// Get monthly earning with id from db
func Get(id int64) (me *MonthlyEarning, err error) {

	me = &MonthlyEarning{
		Id: id,
	}

	o := orm.NewOrm()
	err = o.Read(me)
	if err != nil {
		beego.Error("Failed to read monthly earning:", err)
		return nil, fmt.Errorf("Failed to read monthly earning: %v", err)
	}

	return
}

// Gets all monthly earnings from the database
func GetAll() (mes []*MonthlyEarning, err error) {
	me := MonthlyEarning{}
	o := orm.NewOrm()
	_, err = o.QueryTable(me.TableName()).OrderBy("-Id").All(&mes)

	return
}

// Deletes an monthly earning by ID
func Delete(id int64) error {
	me := MonthlyEarning{
		Id: id,
	}
	o := orm.NewOrm()
	num, err := o.Delete(&me)
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete monthly earning: %v", err))
	}
	beego.Trace("Deleted num monthly earnings:", num)
	return nil
}

func (this *MonthlyEarning) getFileName(interval lib.Interval) string {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 10)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return fmt.Sprintf("invoice-%s-%s", interval.String(), string(b))
}

func (this *MonthlyEarning) getPurchases(locationId int64, interval lib.Interval) (ps []*purchases.Purchase, err error) {
	machines, err := machine.GetAllMachines()
	if err != nil {
		return nil, fmt.Errorf("Failed to get machines: %v", err)
	}
	machinesById := make(map[int64]*machine.Machine)
	for _, machine := range machines {
		machinesById[machine.Id] = machine
	}

	usrs, err := users.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("Failed to get users: %v", err)
	}
	usersById := make(map[int64]users.User)
	for _, user := range usrs {
		usersById[user.Id] = *user
	}

	userMemberships, err := models.GetAllUserMembershipsAt(locationId)
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

	memberships, err := models.GetAllMembershipsAt(locationId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get memberships: %v", err)
	}
	membershipsById := make(map[int64]*models.Membership)
	for _, membership := range memberships {
		membershipsById[membership.Id] = membership
	}

	// Get all uninvoiced purchases in the time range
	ps, err = purchases.GetAllBetweenAt(locationId, interval)
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

func (this *MonthlyEarning) GetInvoices(ps purchases.Purchases) (invs []*Invoice, err error) {

	// Create a slice for unique user summaries.
	users, err := users.GetAllUsers()
	if err != nil {
		return nil, err
	}
	invs = make([]*Invoice, 0, len(users))
	for _, user := range users {
		inv := Invoice{
			User: *user,
		}
		invs = append(invs, &inv)
	}

	// Sort purchases by user.
	for _, p := range ps.Data {

		invExists := false
		var foundInv *Invoice

		for _, inv := range invs {
			if p.User.Id == inv.User.Id {
				invExists = true
				foundInv = inv
				break
			}
		}

		// Create new invoice if it does not exist for the user.
		if !invExists {
			beego.Warn("Creating invoice for purchase that has no matching user")
			newInv := Invoice{
				User: p.User,
			}
			invs = append(invs, &newInv)
			foundInv = invs[len(invs)-1]
		}

		// Append the purchase to the invoice.
		if foundInv.User.Id == p.User.Id {
			foundInv.Purchases.Data = append(foundInv.Purchases.Data, p)
		}
	}

	return
}

func (this *MonthlyEarning) enhancePurchase(purchase *purchases.Purchase,
	machinesById map[int64]*machine.Machine, usersById map[int64]users.User,
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