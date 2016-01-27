package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineTypes_20160127_155843 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineTypes_20160127_155843{}
	m.Created = "20160127_155843"
	migration.Register("MachineTypes_20160127_155843", m)
}

// Run the migrations
func (m *MachineTypes_20160127_155843) Up() {
	m.SQL(`
		CREATE TABLE machine_types (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			short_name varchar(20),
			name varchar(255),
			PRIMARY KEY (id)
	)`)
	m.SQL(`INSERT INTO machine_types VALUES (1, "3dprinter", "3D Printer")`)
	m.SQL(`INSERT INTO machine_types VALUES (2, "cnc", "CNC Mill")`)
	m.SQL(`INSERT INTO machine_types VALUES (3, "heatpress", "Heatpress")`)
	m.SQL(`INSERT INTO machine_types VALUES (4, "knitting", "Knitting Machine")`)
	m.SQL(`INSERT INTO machine_types VALUES (5, "lasercutter", "Lasercutter")`)
	m.SQL(`INSERT INTO machine_types VALUES (6, "vinylcutter", "Vinylcutter")`)
	m.SQL(`ALTER TABLE machines ADD COLUMN type_id int(11) unsigned AFTER type`)
	m.SQL(`ALTER TABLE machines DROP COLUMN type`)
}

// Reverse the migrations
func (m *MachineTypes_20160127_155843) Down() {
	m.SQL(`DROP TABLE machine_types`)
	m.SQL(`ALTER TABLE machines DROP COLUMN type_id`)
	m.SQL(`ALTER TABLE machines ADD COLUMN type varchar(255)`)
}
