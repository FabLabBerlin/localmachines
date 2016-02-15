package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ReservationRulesLocationId_20160215_191406 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ReservationRulesLocationId_20160215_191406{}
	m.Created = "20160215_191406"
	migration.Register("ReservationRulesLocationId_20160215_191406", m)
}

// Run the migrations
func (m *ReservationRulesLocationId_20160215_191406) Up() {
	m.SQL(`
	    ALTER TABLE reservation_rules
	      ADD COLUMN location_id int(11) AFTER id
	`)
	m.SQL(`UPDATE reservation_rules SET location_id = 1`)
}

// Reverse the migrations
func (m *ReservationRulesLocationId_20160215_191406) Down() {
	m.SQL("ALTER TABLE reservation_rules DROP COLUMN location_id")
}
