package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PurchaseTimeEndActual_20151202_102231 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PurchaseTimeEndActual_20151202_102231{}
	m.Created = "20151202_102231"
	migration.Register("PurchaseTimeEndActual_20151202_102231", m)
}

// Run the migrations
func (m *PurchaseTimeEndActual_20151202_102231) Up() {
	m.Sql("ALTER TABLE purchases ADD COLUMN time_end_actual datetime AFTER time_end")
	m.Sql("ALTER TABLE purchases CHANGE COLUMN activation_running running tinyint(1)")
}

// Reverse the migrations
func (m *PurchaseTimeEndActual_20151202_102231) Down() {
	m.Sql("ALTER TABLE purchases DROP COLUMN time_end_actual")
	m.Sql("ALTER TABLE purchases CHANGE COLUMN running activation_running tinyint(1)")
}
