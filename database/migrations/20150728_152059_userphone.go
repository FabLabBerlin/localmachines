package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Userphone_20150728_152059 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Userphone_20150728_152059{}
	m.Created = "20150728_152059"
	migration.Register("Userphone_20150728_152059", m)
}

func (m *Userphone_20150728_152059) Up() {
	m.SQL("ALTER TABLE user ADD COLUMN phone VARCHAR(50)")
}

func (m *Userphone_20150728_152059) Down() {
	m.SQL("ALTER TABLE user DROP COLUMN phone")
}
