package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Reservations_20150928_181415 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Reservations_20150928_181415{}
	m.Created = "20150928_181415"
	migration.Register("Reservations_20150928_181415", m)
}

// Run the migrations
func (m *Reservations_20150928_181415) Up() {
	m.SQL(`CREATE TABLE reservations (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		machine_id int(11) NOT NULL,
		user_id int(11) NOT NULL,
		time_start datetime NOT NULL,
		time_end datetime NOT NULL,
		created datetime NOT NULL,
		PRIMARY KEY (id)
	)`)
}

// Reverse the migrations
func (m *Reservations_20150928_181415) Down() {
	m.SQL("DROP TABLE reservations")
}
