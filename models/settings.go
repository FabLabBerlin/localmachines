package models

import (
	"github.com/astaxie/beego/orm"
)

type Settings struct {
	Data []*Setting
}

func (this *Settings) GetFloat(name string) *float64 {
	if setting := this.getSetting(name); setting != nil {
		return setting.ValueFloat
	}
	return nil
}

func (this *Settings) GetInt(name string) *int64 {
	if setting := this.getSetting(name); setting != nil {
		return setting.ValueInt
	}
	return nil
}

func (this *Settings) GetString(name string) *string {
	if setting := this.getSetting(name); setting != nil {
		return setting.ValueString
	}
	return nil
}

func (this *Settings) getSetting(name string) *Setting {
	for _, setting := range this.Data {
		if setting.Name == name {
			return setting
		}
	}
	return nil
}

type Setting struct {
	Id          int64  `orm:"auto";"pk"`
	Name        string `orm:"size(100)"`
	ValueInt    *int64
	ValueString *string
	ValueFloat  *float64
}

func (this *Setting) TableName() string {
	return "settings"
}

func CreateSetting(setting *Setting) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(setting)
}

func GetAllSettings() (settings Settings, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(Setting).TableName()).All(&settings.Data)
	return
}

func (this *Setting) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(this)
	return
}

func init() {
	orm.RegisterModel(new(Setting))
}
