package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type NetswitchMfi_20160120_174425 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &NetswitchMfi_20160120_174425{}
	m.Created = "20160120_174425"
	migration.Register("NetswitchMfi_20160120_174425", m)
}

// Run the migrations
func (m *NetswitchMfi_20160120_174425) Up() {
	m.Sql("ALTER TABLE netswitch ADD COLUMN host varchar(255) AFTER url_off")
	m.Sql("ALTER TABLE netswitch ADD COLUMN sensor_port int(5) AFTER url_off")
	m.Sql("UPDATE netswitch SET sensor_port = 1")
}

// Reverse the migrations
func (m *NetswitchMfi_20160120_174425) Down() {
	m.Sql("ALTER TABLE netswitch DROP COLUMN host")
	m.Sql("ALTER TABLE netswitch DROP COLUMN sensor_port")
}
