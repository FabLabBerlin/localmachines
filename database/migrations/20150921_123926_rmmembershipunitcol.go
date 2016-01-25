package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Rmmembershipunitcol_20150921_123926 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Rmmembershipunitcol_20150921_123926{}
	m.Created = "20150921_123926"
	migration.Register("Rmmembershipunitcol_20150921_123926", m)
}

// Run the migrations
func (m *Rmmembershipunitcol_20150921_123926) Up() {
	m.SQL("ALTER TABLE membership CHANGE COLUMN duration duration_months INT(11)")
	// Check the unit and if it is days, recalculate it in months
	m.SQL("UPDATE membership SET duration_months=ROUND(duration_months / 30) WHERE unit='days'")
	m.SQL("ALTER TABLE membership DROP COLUMN unit")
}

// Reverse the migrations
func (m *Rmmembershipunitcol_20150921_123926) Down() {
	m.SQL("ALTER TABLE membership ADD COLUMN unit VARCHAR(100) DEFAULT 'months'")
	m.SQL("ALTER TABLE membership CHANGE COLUMN duration_months duration INT(11)")
}
