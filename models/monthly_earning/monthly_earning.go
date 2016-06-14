package monthly_earning

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib"
	"github.com/FabLabBerlin/localmachines/models/monthly_earning/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/settings"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"os"
	"sort"
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
	Invoices    []*invutil.Invoice `orm:"-"`
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

	if !interval.OneMonth() {
		return nil, fmt.Errorf("expected one month")
	}
	m := time.Month(interval.MonthFrom)
	y := interval.YearFrom

	// Create invoices from purchases
	if me.Invoices, err = invutil.GetAllOfMonthAt(locationId, y, m); err != nil {
		err = fmt.Errorf("Failed to get user summaries: %v", err)
		return
	}

	for i := 0; i < len(me.Invoices); i++ {
		sort.Stable(me.Invoices[i].Purchases)
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

func (this *MonthlyEarning) Save() (err error) {
	o := orm.NewOrm()
	this.Id, err = o.Insert(this)
	return
}
