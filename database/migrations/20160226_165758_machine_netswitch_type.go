package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineNetswitchType_20160226_165758 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineNetswitchType_20160226_165758{}
	m.Created = "20160226_165758"
	migration.Register("MachineNetswitchType_20160226_165758", m)
}

// Run the migrations
func (m *MachineNetswitchType_20160226_165758) Up() {
	m.SQL(`ALTER TABLE machines
	  ADD COLUMN netswitch_type VARCHAR(255) AFTER netswitch_sensor_port`)
}

// Reverse the migrations
func (m *MachineNetswitchType_20160226_165758) Down() {
	m.SQL(`ALTER TABLE machines DROP COLUMN netswitch_type`)
}
