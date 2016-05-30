package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PurchasesInvoiceIdAndStatus_20160527_170456 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PurchasesInvoiceIdAndStatus_20160527_170456{}
	m.Created = "20160527_170456"
	migration.Register("PurchasesInvoiceIdAndStatus_20160527_170456", m)
}

// Run the migrations
func (m *PurchasesInvoiceIdAndStatus_20160527_170456) Up() {
	m.SQL("ALTER TABLE purchases ADD COLUMN invoice_id int(11) unsigned")
	m.SQL("ALTER TABLE purchases ADD COLUMN invoice_status varchar(20)")
	m.SQL("ALTER TABLE user_membership ADD COLUMN invoice_id int(11) unsigned")
	m.SQL("ALTER TABLE user_membership ADD COLUMN invoice_status varchar(20)")
}

// Reverse the migrations
func (m *PurchasesInvoiceIdAndStatus_20160527_170456) Down() {
	m.SQL("ALTER TABLE purchases DROP COLUMN invoice_id")
	m.SQL("ALTER TABLE purchases DROP COLUMN invoice_status")
	m.SQL("ALTER TABLE user_membership DROP COLUMN invoice_id")
	m.SQL("ALTER TABLE user_membership DROP COLUMN invoice_status")
}
