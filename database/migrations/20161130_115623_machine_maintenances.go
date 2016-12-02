package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineMaintenances_20161130_115623 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineMaintenances_20161130_115623{}
	m.Created = "20161130_115623"
	migration.Register("MachineMaintenances_20161130_115623", m)
}

// Run the migrations
func (m *MachineMaintenances_20161130_115623) Up() {
	m.SQL(`
CREATE TABLE machine_maintenances (
	id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	machine_id INT(11) UNSIGNED,
	start DATETIME,
	end DATETIME,
	PRIMARY KEY (id)
)
	`)
}

// Reverse the migrations
func (m *MachineMaintenances_20161130_115623) Down() {
	m.SQL("DROP TABLE machine_maintenances")
}
