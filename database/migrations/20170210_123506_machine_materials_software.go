package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MachineMaterialsSoftware_20170210_123506 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MachineMaterialsSoftware_20170210_123506{}
	m.Created = "20170210_123506"
	migration.Register("MachineMaterialsSoftware_20170210_123506", m)
}

// Run the migrations
func (m *MachineMaterialsSoftware_20170210_123506) Up() {
	m.SQL("ALTER TABLE machines ADD COLUMN materials text AFTER links")
	m.SQL("ALTER TABLE machines ADD COLUMN required_software text AFTER materials")
}

// Reverse the migrations
func (m *MachineMaterialsSoftware_20170210_123506) Down() {
	m.SQL("ALTER TABLE machines DROP COLUMN materials")
	m.SQL("ALTER TABLE machines DROP COLUMN required_software")
}
