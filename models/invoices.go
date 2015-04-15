package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Invoice struct {
	Id          int64  `orm:"auto";"pk"`
	Activations string `orm:type(text)`
	XmlFile     string `orm:size(255)`
}

func (this *Invoice) TableName() string {
	return "invoices"
}

func init() {
	orm.RegisterModel(new(Invoice))
}

func CreateInvoice(startTime time.Time,
	endTime time.Time,
	userId int64,
	includeInvoiced bool) (*Invoice, error) {

	beego.Info("Creating invoice...")

	// Load all the matching activations

	// Load activation users

	// Load user memberships

	// Sort activations per user

	// Sort user activations per machine

	// Calculate total price per activation per machine
	// By taking in account the user membership
	return nil, nil
}
