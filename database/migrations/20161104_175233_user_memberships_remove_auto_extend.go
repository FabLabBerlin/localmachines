package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserMembershipsRemoveAutoExtend_20161104_175233 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserMembershipsRemoveAutoExtend_20161104_175233{}
	m.Created = "20161104_175233"
	migration.Register("UserMembershipsRemoveAutoExtend_20161104_175233", m)
}

// Run the migrations
func (m *UserMembershipsRemoveAutoExtend_20161104_175233) Up() {
	m.SQL("ALTER TABLE user_memberships DROP COLUMN auto_extend")
}

// Reverse the migrations
func (m *UserMembershipsRemoveAutoExtend_20161104_175233) Down() {
	m.SQL("ALTER TABLE user_memberships ADD COLUMN auto_extend TINYINT(1)")
	m.SQL("UPDATE user_memberships SET auto_extend = 0")
	m.SQL("UPDATE user_memberships SET auto_extend = 1 WHERE termination_date IS NULL")
	m.SQL("UPDATE user_memberships SET auto_extend = 1 WHERE year(termination_date) > 0")
}
