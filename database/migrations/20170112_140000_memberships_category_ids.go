package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MembershipsPopulateCategoryIds_20170112_140000 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MembershipsPopulateCategoryIds_20170112_140000{}
	m.Created = "20170112_140000"
	migration.Register("MembershipsPopulateCategoryIds_20170112_140000", m)
}

// Run the migrations
func (m *MembershipsPopulateCategoryIds_20170112_140000) Up() {
	m.SQL("ALTER TABLE membership ADD COLUMN affected_categories text")
}

// Reverse the migrations
func (m *MembershipsPopulateCategoryIds_20170112_140000) Down() {
	m.SQL("ALTER TABLE membership DROP COLUMN affected_categories")
}
