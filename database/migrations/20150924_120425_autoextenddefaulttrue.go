package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Autoextenddefaulttrue_20150924_120425 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Autoextenddefaulttrue_20150924_120425{}
	m.Created = "20150924_120425"
	migration.Register("Autoextenddefaulttrue_20150924_120425", m)
}

// Run the migrations
func (m *Autoextenddefaulttrue_20150924_120425) Up() {
	m.SQL("UPDATE membership SET auto_extend = TRUE, auto_extend_duration_months = 1")
	m.SQL("UPDATE user_membership SET auto_extend = TRUE")
}

// Reverse the migrations
func (m *Autoextenddefaulttrue_20150924_120425) Down() {
	m.SQL("UPDATE membership SET auto_extend = NULL, auto_extend_duration_months = NULL")
	m.SQL("UPDATE user_membership SET auto_extend = NULL")
}
