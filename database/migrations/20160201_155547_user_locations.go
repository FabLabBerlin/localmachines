package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserLocations_20160201_155547 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserLocations_20160201_155547{}
	m.Created = "20160201_155547"
	migration.Register("UserLocations_20160201_155547", m)
}

// Run the migrations
func (m *UserLocations_20160201_155547) Up() {
	m.SQL(`
		CREATE TABLE user_locations (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			location_id int(11) unsigned,
			user_id int(11) unsigned,
			user_role varchar(100),
			archived tinyint(1) DEFAULT 0,
			PRIMARY KEY (id)
	)`)
}

// Reverse the migrations
func (m *UserLocations_20160201_155547) Down() {
	m.SQL("DROP TABLE user_locations")
}
