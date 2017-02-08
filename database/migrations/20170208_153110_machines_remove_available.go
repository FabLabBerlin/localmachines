package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachinesRemoveAvailable_20170208_153110 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachinesRemoveAvailable_20170208_153110{}
	m.Created = "20170208_153110"
	migration.Register("MachinesRemoveAvailable_20170208_153110", m)
}

// Run the migrations
func (m *MachinesRemoveAvailable_20170208_153110) Up() {
	m.SQL("ALTER TABLE machines MODIFY available tinyint(1)")
}

// Reverse the migrations
func (m *MachinesRemoveAvailable_20170208_153110) Down() {
	m.SQL("ALTER TABLE machines MODIFY available tinyint(1) NOT NULL")
}
