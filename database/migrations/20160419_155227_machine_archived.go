package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineArchived_20160419_155227 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineArchived_20160419_155227{}
	m.Created = "20160419_155227"
	migration.Register("MachineArchived_20160419_155227", m)
}

// Run the migrations
func (m *MachineArchived_20160419_155227) Up() {
	m.SQL("ALTER TABLE machines ADD COLUMN archived TINYINT(1)")
}

// Reverse the migrations
func (m *MachineArchived_20160419_155227) Down() {
	m.SQL("ALTER TABLE machines DROP coluumn archived")
}
