package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PurchaseCancelledFlag_20151124_153407 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PurchaseCancelledFlag_20151124_153407{}
	m.Created = "20151124_153407"
	migration.Register("PurchaseCancelledFlag_20151124_153407", m)
}

// Run the migrations
func (m *PurchaseCancelledFlag_20151124_153407) Up() {
	m.SQL("ALTER TABLE purchases ADD COLUMN cancelled tinyint(1)")
}

// Reverse the migrations
func (m *PurchaseCancelledFlag_20151124_153407) Down() {
	m.SQL("ALTER TABLE purchases DROP COLUMN cancelled")
}
