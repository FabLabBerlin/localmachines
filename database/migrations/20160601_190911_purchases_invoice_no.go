package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PurchasesInvoiceNo_20160601_190911 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PurchasesInvoiceNo_20160601_190911{}
	m.Created = "20160601_190911"
	migration.Register("PurchasesInvoiceNo_20160601_190911", m)
}

// Run the migrations
func (m *PurchasesInvoiceNo_20160601_190911) Up() {
	m.SQL("ALTER TABLE purchases CHANGE invoice_id invoice_no int(11) unsigned")
	m.SQL("ALTER TABLE user_membership CHANGE invoice_id invoice_no int(11) unsigned")

}

// Reverse the migrations
func (m *PurchasesInvoiceNo_20160601_190911) Down() {
	m.SQL("ALTER TABLE purchases CHANGE invoice_no invoice_id int(11) unsigned")
	m.SQL("ALTER TABLE user_membership CHANGE invoice_no invoice_id int(11) unsigned")
}
