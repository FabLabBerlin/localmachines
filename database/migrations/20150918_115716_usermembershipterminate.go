package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Usermembershipterminate_20150918_115716 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Usermembershipterminate_20150918_115716{}
	m.Created = "20150918_115716"
	migration.Register("Usermembershipterminate_20150918_115716", m)
}

// Run the migrations
func (m *Usermembershipterminate_20150918_115716) Up() {
	m.SQL("ALTER TABLE user_membership ADD COLUMN is_terminated TINYINT(1)")
}

// Reverse the migrations
func (m *Usermembershipterminate_20150918_115716) Down() {
	m.SQL("ALTER TABLE user_membership DROP COLUMN is_terminated")
}
