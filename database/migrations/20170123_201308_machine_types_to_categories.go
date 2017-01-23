package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineTypesToCategories_20170123_201308 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineTypesToCategories_20170123_201308{}
	m.Created = "20170123_201308"
	migration.Register("MachineTypesToCategories_20170123_201308", m)
}

// Run the migrations
func (m *MachineTypesToCategories_20170123_201308) Up() {
	m.SQL("ALTER TABLE machine_types RENAME categories")
}

// Reverse the migrations
func (m *MachineTypesToCategories_20170123_201308) Down() {
	m.SQL("ALTER TABLE categories RENAME machine_types")
}
