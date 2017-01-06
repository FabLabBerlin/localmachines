package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineTypeArchived_20170106_145406 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineTypeArchived_20170106_145406{}
	m.Created = "20170106_145406"
	migration.Register("MachineTypeArchived_20170106_145406", m)
}

// Run the migrations
func (m *MachineTypeArchived_20170106_145406) Up() {
	m.SQL("ALTER TABLE machine_types ADD COLUMN archived TINYINT(1) DEFAULT 0")
}

// Reverse the migrations
func (m *MachineTypeArchived_20170106_145406) Down() {
	m.SQL("ALTER TABLE machine_types DROP COLUMN archived")
}
