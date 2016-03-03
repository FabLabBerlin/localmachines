package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type LocationsLocalIp_20160303_195458 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &LocationsLocalIp_20160303_195458{}
	m.Created = "20160303_195458"
	migration.Register("LocationsLocalIp_20160303_195458", m)
}

// Run the migrations
func (m *LocationsLocalIp_20160303_195458) Up() {
	m.SQL("ALTER TABLE locations ADD COLUMN local_ip VARCHAR(255)")
	m.SQL("UPDATE locations SET local_ip = '37.44.7.170' WHERE id = 1")
}

// Reverse the migrations
func (m *LocationsLocalIp_20160303_195458) Down() {
	m.SQL("ALTER TABLE locations DROP COLUMN local_ip")
}
