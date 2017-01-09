package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineTypeLocationId_20170109_102614 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineTypeLocationId_20170109_102614{}
	m.Created = "20170109_102614"
	migration.Register("MachineTypeLocationId_20170109_102614", m)
}

// Run the migrations
func (m *MachineTypeLocationId_20170109_102614) Up() {
	m.SQL("ALTER TABLE machine_types ADD COLUMN location_id int(11) AFTER id")
	m.SQL(`
INSERT INTO machine_types (location_id, short_name, name, archived)
SELECT l.id, t.short_name, t.name, t.archived
FROM locations AS l,
     machine_types AS t
`)
	m.SQL("DELETE FROM machine_types WHERE location_id IS NULL")
}

// Reverse the migrations
func (m *MachineTypeLocationId_20170109_102614) Down() {
	m.SQL("DELETE FROM machine_types WHERE location_id <> 1")
	m.SQL("ALTER TABLE machine_types DROP COLUMN location_id")
}
