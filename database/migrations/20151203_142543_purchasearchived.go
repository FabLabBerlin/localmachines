package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Purchasearchived_20151203_142543 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Purchasearchived_20151203_142543{}
	m.Created = "20151203_142543"
	migration.Register("Purchasearchived_20151203_142543", m)
}

// Run the migrations
func (m *Purchasearchived_20151203_142543) Up() {
	m.SQL("ALTER TABLE purchases ADD COLUMN archived TINYINT(1) DEFAULT 0")
}

// Reverse the migrations
func (m *Purchasearchived_20151203_142543) Down() {
	m.SQL("ALTER TABLE purchases DROP COLUMN archived")
}
