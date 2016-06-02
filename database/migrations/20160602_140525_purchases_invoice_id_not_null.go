package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PurchasesInvoiceIdNotNull_20160602_140525 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PurchasesInvoiceIdNotNull_20160602_140525{}
	m.Created = "20160602_140525"
	migration.Register("PurchasesInvoiceIdNotNull_20160602_140525", m)
}

// Run the migrations
func (m *PurchasesInvoiceIdNotNull_20160602_140525) Up() {
	m.SQL("ALTER TABLE purchases CHANGE invoice_id invoice_id int(11) unsigned NOT NULL")
	m.SQL("ALTER TABLE user_membership CHANGE invoice_id invoice_id int(11) unsigned NOT NULL")
}

// Reverse the migrations
func (m *PurchasesInvoiceIdNotNull_20160602_140525) Down() {
	m.SQL("ALTER TABLE purchases CHANGE invoice_id invoice_id int(11) unsigned")
	m.SQL("ALTER TABLE user_membership CHANGE invoice_id invoice_id int(11) unsigned")
}
