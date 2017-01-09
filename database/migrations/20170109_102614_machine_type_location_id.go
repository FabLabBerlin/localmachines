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
	m.SQL("ALTER TABLE machine_types ADD COLUMN old_id int(11)")
	m.SQL(`
INSERT INTO machine_types (location_id, short_name, name, archived, old_id)
SELECT l.id, t.short_name, t.name, t.archived, t.id
FROM locations AS l,
     machine_types AS t
`)
	m.SQL("DELETE FROM machine_types WHERE location_id IS NULL")
	m.SQL("UPDATE machines JOIN machine_types AS t ON machines.type_id = old_id SET type_id = t.id")
	m.SQL("UPDATE permission JOIN machine_types AS t ON permission.category_id = old_id SET category_id = t.id")
}

// Reverse the migrations
func (m *MachineTypeLocationId_20170109_102614) Down() {
	m.SQL("UPDATE permission JOIN machine_types AS t ON permission.category_id = t.id SET category_id = old_id")
	m.SQL("UPDATE machines JOIN machine_types AS t ON machines.type_id = t.id SET type_id = old_id")
	m.SQL("DELETE FROM machine_types WHERE location_id <> 1")
	m.SQL("ALTER TABLE machine_types DROP COLUMN location_id")
	m.SQL("UPDATE machine_types SET id = old_id")
	m.SQL("ALTER TABLE machine_types DROP COLUMN old_id")
}
