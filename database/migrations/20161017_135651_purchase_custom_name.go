package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PurchaseCustomName_20161017_135651 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PurchaseCustomName_20161017_135651{}
	m.Created = "20161017_135651"
	migration.Register("PurchaseCustomName_20161017_135651", m)
}

// Run the migrations
func (m *PurchaseCustomName_20161017_135651) Up() {
	m.SQL("ALTER TABLE purchases ADD COLUMN custom_name varchar(255) AFTER machine_id")
}

// Reverse the migrations
func (m *PurchaseCustomName_20161017_135651) Down() {
	m.SQL("ALTER TABLE purchases DROP COLUMN custom_name")
}
