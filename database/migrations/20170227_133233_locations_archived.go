package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type LocationsArchived_20170227_133233 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &LocationsArchived_20170227_133233{}
	m.Created = "20170227_133233"
	migration.Register("LocationsArchived_20170227_133233", m)
}

// Run the migrations
func (m *LocationsArchived_20170227_133233) Up() {
	m.SQL("ALTER TABLE locations ADD COLUMN archived tinyint(1)")
}

// Reverse the migrations
func (m *LocationsArchived_20170227_133233) Down() {
	m.SQL("ALTER TABLE locations DROP COLUMN archived")
}
