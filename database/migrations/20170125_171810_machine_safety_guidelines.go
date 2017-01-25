package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineSafetyGuidelines_20170125_171810 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineSafetyGuidelines_20170125_171810{}
	m.Created = "20170125_171810"
	migration.Register("MachineSafetyGuidelines_20170125_171810", m)
}

// Run the migrations
func (m *MachineSafetyGuidelines_20170125_171810) Up() {
	m.SQL("ALTER TABLE machines ADD COLUMN safety_guidelines text AFTER description")
}

// Reverse the migrations
func (m *MachineSafetyGuidelines_20170125_171810) Down() {
	m.SQL("ALTER TABLE machines DROP COLUMN safety_guidelines")
}
