package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Machinecategory_20150928_125043 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Machinecategory_20150928_125043{}
	m.Created = "20150928_125043"
	migration.Register("Machinecategory_20150928_125043", m)
}

// Run the migrations
func (m *Machinecategory_20150928_125043) Up() {
	m.Sql("ALTER TABLE machines ADD category VARCHAR(30)")

}

// Reverse the migrations
func (m *Machinecategory_20150928_125043) Down() {
	m.Sql("ALTER TABLE machines DROP COLUMN category")

}
