package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type NetswitchRequestType_20151214_175707 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &NetswitchRequestType_20151214_175707{}
	m.Created = "20151214_175707"
	migration.Register("NetswitchRequestType_20151214_175707", m)
}

// Run the migrations
func (m *NetswitchRequestType_20151214_175707) Up() {
	m.Sql("ALTER TABLE netswitch ADD COLUMN request_type varchar(10) NOT NULL DEFAULT 'GET'")
	m.Sql("ALTER TABLE netswitch ADD COLUMN request_data varchar(255)")
}

// Reverse the migrations
func (m *NetswitchRequestType_20151214_175707) Down() {
	m.Sql("ALTER TABLE netswitch DROP COLUMN request_type")
	m.Sql("ALTER TABLE netswitch DROP COLUMN request_data")
}
