package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineType_20160126_151613 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineType_20160126_151613{}
	m.Created = "20160126_151613"
	migration.Register("MachineType_20160126_151613", m)
}

// Run the migrations
func (m *MachineType_20160126_151613) Up() {
	m.SQL("ALTER TABLE machines ADD COLUMN type varchar(255)")
	m.SQL("ALTER TABLE machines ADD COLUMN brand varchar(255)")
	m.SQL("ALTER TABLE machines ADD COLUMN dimensions varchar(255)")
	m.SQL("ALTER TABLE machines ADD COLUMN workspace_dimensions varchar(255)")
}

// Reverse the migrations
func (m *MachineType_20160126_151613) Down() {
	m.SQL("ALTER TABLE machines DROP COLUMN type")
	m.SQL("ALTER TABLE machines DROP COLUMN brand")
	m.SQL("ALTER TABLE machines DROP COLUMN dimensions")
	m.SQL("ALTER TABLE machines DROP COLUMN workspace_dimensions")
}
