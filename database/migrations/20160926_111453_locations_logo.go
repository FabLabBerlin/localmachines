package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type LocationsLogo_20160926_111453 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &LocationsLogo_20160926_111453{}
	m.Created = "20160926_111453"
	migration.Register("LocationsLogo_20160926_111453", m)
}

// Run the migrations
func (m *LocationsLogo_20160926_111453) Up() {
	m.SQL("ALTER TABLE locations ADD COLUMN logo varchar(255)")
}

// Reverse the migrations
func (m *LocationsLogo_20160926_111453) Down() {
	m.SQL("ALTER TABLE locations DROP COLUMN logo")
}
