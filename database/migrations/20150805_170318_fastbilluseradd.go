package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Fastbilluseradd_20150805_170318 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Fastbilluseradd_20150805_170318{}
	m.Created = "20150805_170318"
	migration.Register("Fastbilluseradd_20150805_170318", m)
}

// Run the migrations
func (m *Fastbilluseradd_20150805_170318) Up() {
	m.Sql("ALTER TABLE user ADD COLUMN zip_code VARCHAR(100)")
	m.Sql("ALTER TABLE user ADD COLUMN city VARCHAR(100)")
	m.Sql("ALTER TABLE user ADD COLUMN country_code VARCHAR(2)")
}

// Reverse the migrations
func (m *Fastbilluseradd_20150805_170318) Down() {
	m.Sql("ALTER TABLE user DROP COLUMN zip_code")
	m.Sql("ALTER TABLE user DROP COLUMN city")
	m.Sql("ALTER TABLE user DROP COLUMN country_code")
}
