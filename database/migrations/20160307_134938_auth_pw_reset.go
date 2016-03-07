package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type AuthPwReset_20160307_134938 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AuthPwReset_20160307_134938{}
	m.Created = "20160307_134938"
	migration.Register("AuthPwReset_20160307_134938", m)
}

// Run the migrations
func (m *AuthPwReset_20160307_134938) Up() {
	m.SQL("ALTER TABLE auth ADD COLUMN pw_reset_key VARCHAR(255)")
	m.SQL("ALTER TABLE auth ADD COLUMN pw_reset_time DATETIME")
}

// Reverse the migrations
func (m *AuthPwReset_20160307_134938) Down() {
	m.SQL("ALTER TABLE auth DROP COLUMN pw_reset_key")
	m.SQL("ALTER TABLE auth DROP COLUMN pw_reset_time")
}
