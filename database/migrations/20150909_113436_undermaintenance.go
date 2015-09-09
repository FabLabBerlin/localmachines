package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Undermaintenance_20150909_113436 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Undermaintenance_20150909_113436{}
	m.Created = "20150909_113436"
	migration.Register("Undermaintenance_20150909_113436", m)
}

// Run the migrations
func (m *Undermaintenance_20150909_113436) Up() {
	m.Sql("ALTER TABLE machines ADD under_maintenance tinyint(1)")
}

// Reverse the migrations
func (m *Undermaintenance_20150909_113436) Down() {
	m.Sql("ALTER TABLE machines DROP COLUMN under_maintenance")
}
