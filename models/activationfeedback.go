package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(ActivationFeedback))
}

type ActivationFeedback struct {
	Id           int64 `orm:"auto";"pk"`
	ActivationId int64
	Satisfaction string `orm:"size(100)"`
}

func (mf *ActivationFeedback) TableName() string {
	return "activation_feedback"
}

func CreateActivationFeedback(activationId int64, satisfaction string) (int64, error) {
	o := orm.NewOrm()
	af := ActivationFeedback{
		ActivationId: activationId,
		Satisfaction: satisfaction,
	}
	return o.Insert(&af)
}
