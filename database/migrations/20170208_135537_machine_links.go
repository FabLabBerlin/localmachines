package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineLinks_20170208_135537 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineLinks_20170208_135537{}
	m.Created = "20170208_135537"
	migration.Register("MachineLinks_20170208_135537", m)
}

// Run the migrations
func (m *MachineLinks_20170208_135537) Up() {
	m.SQL("ALTER TABLE machines ADD COLUMN links text AFTER safety_guidelines")
}

// Reverse the migrations
func (m *MachineLinks_20170208_135537) Down() {
	m.SQL("ALTER TABLE machines DROP COLUMN links")
}
