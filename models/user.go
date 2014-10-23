package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(User))
}

type User struct {
	Id         int    `orm:"auto";"pk"`
	FirstName  string `orm:"size(100)"`
	LastName   string `orm:"size(100)"`
	Username   string `orm:"size(100)"`
	Email      string `orm:"size(100)"`
	IvoiceAddr int
	ShipAddr   int
	ClientId   int
	B2b        bool
	Company    string `orm:"size(100)"`
	VatUserId  string `orm:"size(100)"`
	VatRate    int
}

func (this *User) GetUserIdForUsername(username string) int {
	o := orm.NewOrm()
	err := o.Raw("SELECT id FROM user WHERE username = ?", username).QueryRow(this)
	if err != nil {
		beego.Error(err)
	}
	return this.Id
}

func (this *User) GetPasswordForUsername(username string) string {
	o := orm.NewOrm()
	data := struct{ Password string }{}
	err := o.Raw("SELECT password FROM auth INNER JOIN user ON auth.user_id = user.id WHERE user.username = ?", username).QueryRow(&data)
	if err != nil {
		beego.Error(err)
	}
	return data.Password
}
