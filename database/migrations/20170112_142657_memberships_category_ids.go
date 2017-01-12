package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MembershipsCategoryIds_20170112_142657 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MembershipsCategoryIds_20170112_142657{}
	m.Created = "20170112_142657"
	migration.Register("MembershipsCategoryIds_20170112_142657", m)
}

// Run the migrations
func (m *MembershipsCategoryIds_20170112_142657) Up() {
	m.SQL("ALTER TABLE membership ADD COLUMN affected_categories text AFTER affected_machines")
}

// Reverse the migrations
func (m *MembershipsCategoryIds_20170112_142657) Down() {
	m.SQL("ALTER TABLE membership DROP COLUMN affected_categories")
}
