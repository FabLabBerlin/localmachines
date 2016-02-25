package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineGracePeriod_20160225_141536 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineGracePeriod_20160225_141536{}
	m.Created = "20160225_141536"
	migration.Register("MachineGracePeriod_20160225_141536", m)
}

// Run the migrations
func (m *MachineGracePeriod_20160225_141536) Up() {
	m.SQL("ALTER TABLE machines ADD COLUMN grace_period int(11) AFTER reservation_price_hourly")
}

// Reverse the migrations
func (m *MachineGracePeriod_20160225_141536) Down() {
	m.SQL("ALTER TABLE machines DROP COLUMN grace_period")
}
