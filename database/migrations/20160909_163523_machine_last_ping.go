package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineLastPing_20160909_163523 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineLastPing_20160909_163523{}
	m.Created = "20160909_163523"
	migration.Register("MachineLastPing_20160909_163523", m)
}

// Run the migrations
func (m *MachineLastPing_20160909_163523) Up() {
	m.SQL("ALTER TABLE machines ADD COLUMN netswitch_last_ping DATETIME")
}

// Reverse the migrations
func (m *MachineLastPing_20160909_163523) Down() {
	m.SQL("ALTER TABLE machines DROP COLUMN netswitch_last_ping")
}
