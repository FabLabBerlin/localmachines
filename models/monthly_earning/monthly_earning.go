package monthly_earning

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/settings"
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
	Id          int64
	LocationId  int64
	MonthFrom   int
	YearFrom    int
	MonthTo     int
	YearTo      int
	Activations string `orm:"type(text)"`
	FilePath    string `orm:"size(255)"`
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
func New(locationId int64, interval lib.Interval) (me *MonthlyEarning, err error) {
	me = &MonthlyEarning{
		LocationId: locationId,
		MonthFrom:  interval.MonthFrom,
		YearFrom:   interval.YearFrom,
		MonthTo:    interval.MonthTo,
		YearTo:     interval.YearTo,
	}

	locSettings, err := settings.GetAllAt(locationId)
	if err != nil {
		return nil, fmt.Errorf("get settings: %v", err)
	}
	var vatPercent float64
	if vat := locSettings.GetFloat(locationId, settings.VAT); vat != nil {
		vatPercent = *vat
	} else {
		vatPercent = 19.0
	}
	beego.Info("vatPercent=", vatPercent)

	// Create invoices from purchases
	if me.Invoices, err = me.NewInvoices(vatPercent); err != nil {
		err = fmt.Errorf("Failed to get user summaries: %v", err)
		return
	}

	for i := 0; i < len(me.Invoices); i++ {
		sort.Stable(me.Invoices[i].Purchases)
		for _, purchase := range me.Invoices[i].Purchases.Data {
			purchase.TotalPrice = purchases.PriceTotalExclDisc(purchase)
			purchase.DiscountedTotal, err = purchases.PriceTotalDisc(purchase)
			if err != nil {
				return nil, fmt.Errorf("price total disc (purchase %v): %v", purchase.Id, err)
			}
		}
	}

	return me, err
}

// Creates monthly earning entry in the database
func Create(locationId int64, interval lib.Interval) (me *MonthlyEarning, err error) {

	if me, err = New(locationId, interval); err != nil {
		return nil, fmt.Errorf("New: %v", err)
	}

	// Create *.xlsx file.
	fileName := me.getFileName(interval)

	dirName := fmt.Sprintf("excel_exports/%v", locationId)

	// The files directory must exist and be writable.
	if exists, _ := exists(dirName); !exists {
		if err = os.MkdirAll(dirName, 0777); err != nil {
			beego.Error("Failed to create excel_exports dir:", err)
			return nil, fmt.Errorf("Failed to create excel_exports dir: %v", err)
		}
	}

	me.FilePath = dirName + "/" + fileName + ".xlsx"

	err = createXlsxFile(me.FilePath, me)
	if err != nil {
		return nil, fmt.Errorf("Failed to create *.xlsx file: %v", err)
	}

	me.Created = time.Now()
	me.MonthFrom = interval.MonthFrom
	me.YearFrom = interval.YearFrom
	me.MonthTo = interval.MonthTo
	me.YearTo = interval.YearTo

	if err = me.Save(); err != nil {
		return nil, fmt.Errorf("save: %v", err)
	}

	return
}

// Get monthly earning with id from db
func Get(id int64) (me *MonthlyEarning, err error) {
	me = &MonthlyEarning{
		Id: id,
	}

	o := orm.NewOrm()
	err = o.Read(me)

	return
}

// Gets all monthly earnings from the database
func GetAllAt(locationId int64) (mes []*MonthlyEarning, err error) {
	me := MonthlyEarning{}
	o := orm.NewOrm()
	_, err = o.QueryTable(me.TableName()).
		OrderBy("-Id").
		Filter("location_id", locationId).
		All(&mes)

	return
}

// Deletes a monthly earning by Id
func Delete(id int64) (err error) {
	me := MonthlyEarning{
		Id: id,
	}
	o := orm.NewOrm()
	_, err = o.Delete(&me)
	return
}

func (this *MonthlyEarning) getFileName(interval lib.Interval) string {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 10)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return fmt.Sprintf("monthly-earnings-%s-%s", interval.String(), string(b))
}

func (this *MonthlyEarning) getPurchases(locationId int64, interval lib.Interval) (ps []*purchases.Purchase, err error) {
	machines, err := machine.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to get machines: %v", err)
	}
	machinesById := make(map[int64]*machine.Machine)
	for _, machine := range machines {
		machinesById[machine.Id] = machine
	}

	usrs, err := users.GetAllUsersAt(locationId)
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

func (this *MonthlyEarning) NewInvoices(vatPercent float64) (invs []*Invoice, err error) {
	// Enhance activations with user and membership data
	ps, err := this.getPurchases(this.LocationId, this.Interval())
	if err != nil {
		return nil, fmt.Errorf("Failed to get enhanced purchaes: %v", err)
	}

	activationIds := make([]string, 0, len(ps))
	for _, p := range ps {
		if p.Type == purchases.TYPE_ACTIVATION {
			activationIds = append(activationIds, strconv.FormatInt(p.Id, 10))
		}
	}
	this.Activations = "[" + strings.Join(activationIds, ",") + "]"

	// Create a slice for unique user summaries.
	users, err := users.GetAllUsersAt(this.LocationId)
	if err != nil {
		return nil, err
	}
	invs = make([]*Invoice, 0, len(users))
	for _, user := range users {
		inv := Invoice{
			User:       *user,
			VatPercent: vatPercent,
		}
		invs = append(invs, &inv)
	}

	// Sort purchases by user.
	for _, p := range ps {

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
		return fmt.Errorf("No user has the ID %v", purchase.UserId)
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

		if usrMem.Interval().Contains(purchase.TimeStart) {
			purchase.Memberships = append(purchase.Memberships, mem)
		}
	}

	return nil
}

func (this *MonthlyEarning) Save() (err error) {
	o := orm.NewOrm()
	this.Id, err = o.Insert(this)
	return
}
