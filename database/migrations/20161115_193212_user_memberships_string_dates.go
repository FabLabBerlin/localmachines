package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserMembershipsStringDates_20161115_193212 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserMembershipsStringDates_20161115_193212{}
	m.Created = "20161115_193212"
	migration.Register("UserMembershipsStringDates_20161115_193212", m)
}

// Run the migrations
func (m *UserMembershipsStringDates_20161115_193212) Up() {
	m.SQL("ALTER TABLE user_memberships CHANGE start_date start_date VARCHAR(255)")
	m.SQL("UPDATE user_memberships SET start_date = SUBSTR(start_date, 1, 10)")
	m.SQL("ALTER TABLE user_memberships CHANGE termination_date termination_date VARCHAR(255)")
	m.SQL("UPDATE user_memberships SET termination_date = SUBSTR(termination_date, 1, 10)")

	m.SQL("ALTER TABLE invoice_user_memberships CHANGE start_date start_date VARCHAR(255)")
	m.SQL("UPDATE invoice_user_memberships SET start_date = SUBSTR(start_date, 1, 10)")
	m.SQL("ALTER TABLE invoice_user_memberships CHANGE termination_date termination_date VARCHAR(255)")
	m.SQL("UPDATE invoice_user_memberships SET termination_date = SUBSTR(termination_date, 1, 10)")
}

// Reverse the migrations
func (m *UserMembershipsStringDates_20161115_193212) Down() {
	m.SQL("ALTER TABLE user_memberships CHANGE start_date start_date DATETIME")
	m.SQL("ALTER TABLE user_memberships CHANGE termination_date termination_date DATETIME")
	m.SQL("ALTER TABLE invoice_user_memberships CHANGE start_date start_date DATETIME")
	m.SQL("ALTER TABLE invoice_user_memberships CHANGE termination_date termination_date DATETIME")
}
