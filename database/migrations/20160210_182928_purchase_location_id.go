package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PurchaseLocationId_20160210_182928 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PurchaseLocationId_20160210_182928{}
	m.Created = "20160210_182928"
	migration.Register("PurchaseLocationId_20160210_182928", m)
}

// Run the migrations
func (m *PurchaseLocationId_20160210_182928) Up() {
	m.SQL(`
	    ALTER TABLE purchases 
	      ADD COLUMN location_id INT(11) after id
	`)
	m.SQL(`UPDATE purchases SET location_id = 1`)
}

// Reverse the migrations
func (m *PurchaseLocationId_20160210_182928) Down() {
	m.SQL("ALTER TABLE purchases DROP COLUMN location_id")
}
