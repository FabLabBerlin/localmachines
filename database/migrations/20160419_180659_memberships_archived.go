package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MembershipsArchived_20160419_180659 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MembershipsArchived_20160419_180659{}
	m.Created = "20160419_180659"
	migration.Register("MembershipsArchived_20160419_180659", m)
}

// Run the migrations
func (m *MembershipsArchived_20160419_180659) Up() {
	m.SQL("ALTER TABLE membership ADD COLUMN archived TINYINT(1)")
}

// Reverse the migrations
func (m *MembershipsArchived_20160419_180659) Down() {
	m.SQL("ALTER TABLE membership DROP coluumn archived")
}
