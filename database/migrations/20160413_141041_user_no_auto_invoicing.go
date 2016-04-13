package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserNoAutoInvoicing_20160413_141041 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserNoAutoInvoicing_20160413_141041{}
	m.Created = "20160413_141041"
	migration.Register("UserNoAutoInvoicing_20160413_141041", m)
}

// Run the migrations
func (m *UserNoAutoInvoicing_20160413_141041) Up() {
	m.SQL("ALTER TABLE user ADD COLUMN no_auto_invoicing TINYINT(1)")

}

// Reverse the migrations
func (m *UserNoAutoInvoicing_20160413_141041) Down() {
	m.SQL("ALTER TABLE user DROP COLUMN no_auto_invoicing")
}
