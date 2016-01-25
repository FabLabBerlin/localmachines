package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Billingingaddress_20150728_120222 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Billingingaddress_20150728_120222{}
	m.Created = "20150728_120222"
	migration.Register("Billingingaddress_20150728_120222", m)
}

// Run the migrations
func (m *Billingingaddress_20150728_120222) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE user MODIFY invoice_addr TEXT")
	m.SQL("ALTER TABLE user MODIFY ship_addr TEXT")
}

// Reverse the migrations
func (m *Billingingaddress_20150728_120222) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE user MODIFY invoice_addr INT(11)")
	m.SQL("ALTER TABLE user MODIFY ship_addr INT(11)")
}
