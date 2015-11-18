package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Reservationdisabled_20151118_133537 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Reservationdisabled_20151118_133537{}
	m.Created = "20151118_133537"
	migration.Register("Reservationdisabled_20151118_133537", m)
}

// Run the migrations
func (m *Reservationdisabled_20151118_133537) Up() {
	m.Sql("ALTER TABLE reservations ADD COLUMN disabled tinyint(1)")

}

// Reverse the migrations
func (m *Reservationdisabled_20151118_133537) Down() {
	m.Sql("ALTER TABLE reservations DROP COLUMN disabled")
}
