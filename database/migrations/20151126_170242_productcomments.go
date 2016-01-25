package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Productcomments_20151126_170242 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Productcomments_20151126_170242{}
	m.Created = "20151126_170242"
	migration.Register("Productcomments_20151126_170242", m)
}

// Run the migrations
func (m *Productcomments_20151126_170242) Up() {
	m.SQL("ALTER TABLE products ADD COLUMN comments TEXT")
}

// Reverse the migrations
func (m *Productcomments_20151126_170242) Down() {
	m.SQL("ALTER TABLE products DROP COLUMN comments")
}
