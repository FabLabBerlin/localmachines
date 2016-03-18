package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MonthlyEarningLocationId_20160318_170010 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MonthlyEarningLocationId_20160318_170010{}
	m.Created = "20160318_170010"
	migration.Register("MonthlyEarningLocationId_20160318_170010", m)
}

// Run the migrations
func (m *MonthlyEarningLocationId_20160318_170010) Up() {
	m.SQL(`
	    ALTER TABLE monthly_earnings 
	      ADD COLUMN location_id INT(11) after id
	`)
	m.SQL(`UPDATE monthly_earnings SET location_id = 1`)
}

// Reverse the migrations
func (m *MonthlyEarningLocationId_20160318_170010) Down() {
	m.SQL("ALTER TABLE monthly_earnings DROP COLUMN location_id")
}
