package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PurchasesInvoiceStatusNotNull_20160601_153246 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PurchasesInvoiceStatusNotNull_20160601_153246{}
	m.Created = "20160601_153246"
	migration.Register("PurchasesInvoiceStatusNotNull_20160601_153246", m)
}

// Run the migrations
func (m *PurchasesInvoiceStatusNotNull_20160601_153246) Up() {
	m.SQL("ALTER TABLE purchases DROP COLUMN invoice_status")
	m.SQL("ALTER TABLE user_membership DROP COLUMN invoice_status")
	m.SQL("ALTER TABLE purchases ADD COLUMN invoice_status varchar(20) NOT NULL DEFAULT ''")
	m.SQL("ALTER TABLE user_membership ADD COLUMN invoice_status varchar(20) NOT NULL DEFAULT ''")
}

// Reverse the migrations
func (m *PurchasesInvoiceStatusNotNull_20160601_153246) Down() {
	m.SQL("ALTER TABLE purchases DROP COLUMN invoice_status")
	m.SQL("ALTER TABLE user_membership DROP COLUMN invoice_status")
	m.SQL("ALTER TABLE purchases ADD COLUMN invoice_status varchar(20)")
	m.SQL("ALTER TABLE user_membership ADD COLUMN invoice_status varchar(20)")
}
