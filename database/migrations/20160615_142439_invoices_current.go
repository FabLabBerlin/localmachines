package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type InvoicesCurrent_20160615_142439 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &InvoicesCurrent_20160615_142439{}
	m.Created = "20160615_142439"
	migration.Register("InvoicesCurrent_20160615_142439", m)
}

// Run the migrations
func (m *InvoicesCurrent_20160615_142439) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE invoices ADD COLUMN current TINYINT(1) NOT NULL")
}

// Reverse the migrations
func (m *InvoicesCurrent_20160615_142439) Down() {
	m.SQL("ALTER TABLE invoices DROP COLUMN current")
}
