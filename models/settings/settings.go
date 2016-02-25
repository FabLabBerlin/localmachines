package settings

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

const TABLE_NAME = "settings"

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

type Setting struct {
	Id          int64 `orm:"auto";"pk"`
	LocationId  int64
	Name        string `orm:"size(100)"`
	ValueInt    *int64
	ValueString *string
	ValueFloat  *float64
}

func (this *Setting) TableName() string {
	return TABLE_NAME
}

func Create(setting *Setting) (int64, error) {
	o := orm.NewOrm()
	if setting.LocationId <= 0 {
		return 0, fmt.Errorf("location id must be defined")
	}
	return o.Insert(setting)
}

func GetAllAt(locationId int64) (settings Settings, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).
		Filter("location_id", locationId).
		All(&settings.Data)
	return
}

func (this *Setting) Update() (err error) {
	o := orm.NewOrm()
	if this.LocationId <= 0 {
		return fmt.Errorf("location id must be defined")
	}
	_, err = o.Update(this)
	return
}

func init() {
	orm.RegisterModel(new(Setting))
}
