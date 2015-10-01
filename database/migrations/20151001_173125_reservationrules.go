package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Reservationrules_20151001_173125 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Reservationrules_20151001_173125{}
	m.Created = "20151001_173125"
	migration.Register("Reservationrules_20151001_173125", m)
}

// Run the migrations
func (m *Reservationrules_20151001_173125) Up() {
	m.Sql(`CREATE TABLE reservation_rules (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		machine_id int(11),
		available tinyint(1),
		unavailable tinyint(1),
		time_start datetime,
		time_end datetime,
		monday tinyint(1),
		tuesday tinyint(1),
		wednesday tinyint(1),
		thursday tinyint(1),
		friday tinyint(1),
		saturday tinyint(1),
		sunday tinyint(1),
		created datetime NOT NULL,
		PRIMARY KEY (id)
	)`)
}

// Reverse the migrations
func (m *Reservationrules_20151001_173125) Down() {
	m.Sql("DROP TABLE reservation_rules")
}
