package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserUniversity_20161208_132309 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserUniversity_20161208_132309{}
	m.Created = "20161208_132309"
	migration.Register("UserUniversity_20161208_132309", m)
}

// Run the migrations
func (m *UserUniversity_20161208_132309) Up() {
	m.SQL("ALTER TABLE user ADD COLUMN student_id varchar(50)")
	m.SQL("ALTER TABLE user ADD COLUMN security_briefing varchar(255)")
}

// Reverse the migrations
func (m *UserUniversity_20161208_132309) Down() {
	m.SQL("ALTER TABLE user DROP COLUMN student_id")
	m.SQL("ALTER TABLE user DROP COLUMN security_briefing")
}
