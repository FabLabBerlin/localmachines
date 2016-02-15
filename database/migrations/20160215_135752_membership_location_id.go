package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MembershipLocationId_20160215_135752 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MembershipLocationId_20160215_135752{}
	m.Created = "20160215_135752"
	migration.Register("MembershipLocationId_20160215_135752", m)
}

// Run the migrations
func (m *MembershipLocationId_20160215_135752) Up() {
	m.SQL(`
	    ALTER TABLE membership
	      ADD COLUMN location_id int(11) AFTER id
	`)
	m.SQL(`UPDATE membership SET location_id = 1`)
}

// Reverse the migrations
func (m *MembershipLocationId_20160215_135752) Down() {
	m.SQL("ALTER TABLE membership DROP COLUMN location_id")
}
