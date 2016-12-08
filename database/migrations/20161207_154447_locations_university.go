package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type LocationsUniversity_20161207_154447 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &LocationsUniversity_20161207_154447{}
	m.Created = "20161207_154447"
	migration.Register("LocationsUniversity_20161207_154447", m)
}

// Run the migrations
func (m *LocationsUniversity_20161207_154447) Up() {
	m.SQL("ALTER TABLE locations ADD COLUMN university TINYINT(1)")
}

// Reverse the migrations
func (m *LocationsUniversity_20161207_154447) Down() {
	m.SQL("ALTER TABLE locations DROP COLUMN university")
}
