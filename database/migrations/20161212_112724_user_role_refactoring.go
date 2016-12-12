package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserRoleRefactoring_20161212_112724 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserRoleRefactoring_20161212_112724{}
	m.Created = "20161212_112724"
	migration.Register("UserRoleRefactoring_20161212_112724", m)
}

// Run the migrations
func (m *UserRoleRefactoring_20161212_112724) Up() {
	m.SQL("ALTER TABLE user ADD COLUMN super_admin TINYINT(1) AFTER user_role")
	m.SQL("ALTER TABLE user DROP COLUMN user_role")
}

// Reverse the migrations
func (m *UserRoleRefactoring_20161212_112724) Down() {
	m.SQL("ALTER TABLE user ADD COLUMN user_role VARCHAR(20) AFTER super_admin")
	m.SQL("ALTER TABLE user DROP COLUMN super_admin")
}
