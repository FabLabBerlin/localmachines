package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type SettingsLocationId_20160225_163048 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &SettingsLocationId_20160225_163048{}
	m.Created = "20160225_163048"
	migration.Register("SettingsLocationId_20160225_163048", m)
}

// Run the migrations
func (m *SettingsLocationId_20160225_163048) Up() {
	m.SQL(`
	    ALTER TABLE settings 
	      ADD COLUMN location_id INT(11) after id
	`)
	m.SQL(`UPDATE settings SET location_id = 1`)
}

// Reverse the migrations
func (m *SettingsLocationId_20160225_163048) Down() {
	m.SQL("ALTER TABLE settings DROP COLUMN location_id")
}
