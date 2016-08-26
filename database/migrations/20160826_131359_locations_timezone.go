package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type LocationsTimezone_20160826_131359 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &LocationsTimezone_20160826_131359{}
	m.Created = "20160826_131359"
	migration.Register("LocationsTimezone_20160826_131359", m)
}

// Run the migrations
func (m *LocationsTimezone_20160826_131359) Up() {
	m.SQL("ALTER TABLE locations ADD COLUMN timezone varchar(100)")

}

// Reverse the migrations
func (m *LocationsTimezone_20160826_131359) Down() {
	m.SQL("ALTER TABLE locations DROP COLUMN timezone")
}
