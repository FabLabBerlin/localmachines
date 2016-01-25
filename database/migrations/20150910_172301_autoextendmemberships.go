package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Autoextendmemberships_20150910_172301 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Autoextendmemberships_20150910_172301{}
	m.Created = "20150910_172301"
	migration.Register("Autoextendmemberships_20150910_172301", m)
}

// Run the migrations
func (m *Autoextendmemberships_20150910_172301) Up() {
	m.SQL("ALTER TABLE membership ADD COLUMN auto_extend TINYINT(1)")
	m.SQL("ALTER TABLE membership ADD COLUMN auto_extend_duration INT(11)")
}

// Reverse the migrations
func (m *Autoextendmemberships_20150910_172301) Down() {
	m.SQL("ALTER TABLE membership DROP COLUMN auto_extend")
	m.SQL("ALTER TABLE membership DROP COLUMN auto_extend_duration")
}
