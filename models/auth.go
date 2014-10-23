package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	// Register models
	orm.RegisterModel(new(Auth))
}

type Auth struct {
	UserId   int    `orm:"auto"`
	NfcKey   string `orm:"size(100)"`
	Password string `orm:"size(100)"`
}

func (this *Auth) GetPassword(userId int) string {
	o := orm.NewOrm()
	err := o.Raw("SELECT password FROM auth WHERE user_id = ?", userId).QueryRow(this)
	if err != nil {
		beego.Error(err)
	}
	return this.Password
}
