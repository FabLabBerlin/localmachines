package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Usermembershipautoextend_20150910_110015 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Usermembershipautoextend_20150910_110015{}
	m.Created = "20150910_110015"
	migration.Register("Usermembershipautoextend_20150910_110015", m)
}

// Run the migrations
func (m *Usermembershipautoextend_20150910_110015) Up() {
	m.SQL("ALTER TABLE user_membership ADD COLUMN auto_extend TINYINT(1)")
}

// Reverse the migrations
func (m *Usermembershipautoextend_20150910_110015) Down() {
	m.SQL("ALTER TABLE user_membership DROP COLUMN auto_extend")
}
