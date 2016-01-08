package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineXmpp_20160108_113817 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineXmpp_20160108_113817{}
	m.Created = "20160108_113817"
	migration.Register("MachineXmpp_20160108_113817", m)
}

// Run the migrations
func (m *MachineXmpp_20160108_113817) Up() {
	m.Sql("ALTER TABLE machines ADD COLUMN xmpp tinyint(1)")
}

// Reverse the migrations
func (m *MachineXmpp_20160108_113817) Down() {
	m.Sql("ALTER TABLE machines DROP COLUMN xmpp")
}
