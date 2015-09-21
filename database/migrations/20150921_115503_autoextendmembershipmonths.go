package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Autoextendmembershipmonths_20150921_115503 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Autoextendmembershipmonths_20150921_115503{}
	m.Created = "20150921_115503"
	migration.Register("Autoextendmembershipmonths_20150921_115503", m)
}

// Run the migrations
func (m *Autoextendmembershipmonths_20150921_115503) Up() {
	m.Sql("ALTER TABLE membership CHANGE COLUMN " +
		"auto_extend_duration auto_extend_duration_months INT(11)")
}

// Reverse the migrations
func (m *Autoextendmembershipmonths_20150921_115503) Down() {
	m.Sql("ALTER TABLE membership CHANGE COLUMN " +
		"auto_extend_duration_months auto_extend_duration INT(11)")
}
