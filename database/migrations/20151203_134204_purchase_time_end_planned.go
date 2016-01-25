package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PurchaseTimeEndPlanned_20151203_134204 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PurchaseTimeEndPlanned_20151203_134204{}
	m.Created = "20151203_134204"
	migration.Register("PurchaseTimeEndPlanned_20151203_134204", m)
}

// Run the migrations
func (m *PurchaseTimeEndPlanned_20151203_134204) Up() {
	m.SQL("ALTER TABLE purchases CHANGE COLUMN time_end_actual time_end_planned datetime")
}

// Reverse the migrations
func (m *PurchaseTimeEndPlanned_20151203_134204) Down() {
	m.SQL("ALTER TABLE purchases CHANGE COLUMN time_end_planned time_end_actual datetime")
}
