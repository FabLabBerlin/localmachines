package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineImageSmall_20160927_152202 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineImageSmall_20160927_152202{}
	m.Created = "20160927_152202"
	migration.Register("MachineImageSmall_20160927_152202", m)
}

// Run the migrations
func (m *MachineImageSmall_20160927_152202) Up() {
	m.SQL("ALTER TABLE machines ADD COLUMN image_small varchar(255) AFTER image")
	m.SQL("UPDATE machines SET image_small = image")
}

// Reverse the migrations
func (m *MachineImageSmall_20160927_152202) Down() {
	m.SQL("ALTER TABLE machines DROP COLUMN image_small")
}
