package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Locations_20160119_152047 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Locations_20160119_152047{}
	m.Created = "20160119_152047"
	migration.Register("Locations_20160119_152047", m)
}

// Run the migrations
func (m *Locations_20160119_152047) Up() {
	m.SQL(`CREATE TABLE locations (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(100),
		PRIMARY KEY (id)
	)`)
	m.SQL(`INSERT INTO locations VALUES (1, "Fab Lab Berlin")`)
	m.SQL(`ALTER TABLE machines ADD COLUMN location_id int(11) unsigned AFTER id`)
	m.SQL(`ALTER TABLE netswitch ADD COLUMN location_id int(11) unsigned AFTER id`)
	m.SQL(`UPDATE machines SET location_id = 1`)
}

// Reverse the migrations
func (m *Locations_20160119_152047) Down() {
	m.SQL("DROP TABLE locations")
	m.SQL(`ALTER TABLE machines DROP column location_id`)
}
