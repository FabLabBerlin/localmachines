package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Purchasecomments_20151203_160933 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Purchasecomments_20151203_160933{}
	m.Created = "20151203_160933"
	migration.Register("Purchasecomments_20151203_160933", m)
}

// Run the migrations
func (m *Purchasecomments_20151203_160933) Up() {
	m.Sql("ALTER TABLE purchases ADD COLUMN comments TEXT DEFAULT ''")
}

// Reverse the migrations
func (m *Purchasecomments_20151203_160933) Down() {
	m.Sql("ALTER TABLE purchases DROP COLUMN comments")
}
