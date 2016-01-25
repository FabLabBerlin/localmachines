package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Productarchived_20151202_172646 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Productarchived_20151202_172646{}
	m.Created = "20151202_172646"
	migration.Register("Productarchived_20151202_172646", m)
}

// Run the migrations
func (m *Productarchived_20151202_172646) Up() {
	m.SQL("ALTER TABLE products ADD COLUMN archived TINYINT(1) DEFAULT 0")
}

// Reverse the migrations
func (m *Productarchived_20151202_172646) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE products DROP COLUMN archived")
}
