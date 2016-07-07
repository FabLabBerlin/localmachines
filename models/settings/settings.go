package settings

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

const TABLE_NAME = "settings"

const (
	CURRENCY             = "Currency"
	TERMS_URL            = "TermsUrl"
	VAT                  = "VAT"
	FASTBILL_TEMPLATE_ID = "FastbillTemplateId"
)

var validNames = []string{
	CURRENCY,
	TERMS_URL,
	VAT,
	FASTBILL_TEMPLATE_ID,
}

type Settings struct {
	Data []*Setting
}

func (this *Settings) GetFloat(locationId int64, name string) *float64 {
	if setting := this.getSetting(locationId, name); setting != nil {
		return setting.ValueFloat
	}
	return nil
}

func (this *Settings) GetInt(locationId int64, name string) *int64 {
	if setting := this.getSetting(locationId, name); setting != nil {
		return setting.ValueInt
	}
	return nil
}

func (this *Settings) GetString(locationId int64, name string) *string {
	if setting := this.getSetting(locationId, name); setting != nil {
		return setting.ValueString
	}
	return nil
}

func (this *Settings) getSetting(locationId int64, name string) *Setting {
	for _, setting := range this.Data {
		if setting.Name == name {
			return setting
		}
	}
	return nil
}

func isValidName(name string) bool {
	for _, validName := range validNames {
		if name == validName {
			return true
		}
	}
	return false
}

type Setting struct {
	Id          int64
	LocationId  int64
	Name        string `orm:"size(100)"`
	ValueInt    *int64
	ValueString *string
	ValueFloat  *float64
}

func (this *Setting) TableName() string {
	return TABLE_NAME
}

func Create(setting *Setting) (id int64, err error) {
	o := orm.NewOrm()

	if err = o.Begin(); err != nil {
		return 0, fmt.Errorf("begin tx: %v", err)
	}

	if !isValidName(setting.Name) {
		return 0, fmt.Errorf("'%v' is not a valid name", setting.Name)
	}
	if setting.LocationId <= 0 {
		return 0, fmt.Errorf("location id must be defined")
	}
	query := "DELETE FROM settings WHERE location_id = ? AND name = ?"
	_, err = o.Raw(query, setting.LocationId, setting.Name).Exec()
	if err != nil {
		o.Rollback()
		return 0, fmt.Errorf("delete old: %v", err)
	}
	if _, err := o.Insert(setting); err != nil {
		o.Rollback()
		return 0, fmt.Errorf("insert: %v", err)
	}
	if err := o.Commit(); err != nil {
		return 0, fmt.Errorf("commit tx: %v", err)
	}
	return
}

func GetAllAt(locationId int64) (settings Settings, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).
		Filter("location_id", locationId).
		All(&settings.Data)
	return
}

func (this *Setting) Update() (err error) {
	_, err = Create(this)
	return
}

func init() {
	orm.RegisterModel(new(Setting))
}
