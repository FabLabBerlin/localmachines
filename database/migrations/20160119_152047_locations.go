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
	m.Sql(`CREATE TABLE locations (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(100),
		PRIMARY KEY (id)
	)`)
	m.Sql(`INSERT INTO locations VALUES (1, "Fab Lab Berlin")`)
	m.Sql(`ALTER TABLE machines ADD COLUMN location_id int(11) unsigned AFTER id`)
	m.Sql(`ALTER TABLE netswitch ADD COLUMN location_id int(11) unsigned AFTER id`)
	m.Sql(`UPDATE machines SET location_id = 1`)
}

// Reverse the migrations
func (m *Locations_20160119_152047) Down() {
	m.Sql("DROP TABLE locations")
	m.Sql(`ALTER TABLE machines DROP column location_id`)
}
