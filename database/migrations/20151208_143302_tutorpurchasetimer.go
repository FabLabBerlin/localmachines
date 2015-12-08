package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Tutorpurchasetimer_20151208_143302 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Tutorpurchasetimer_20151208_143302{}
	m.Created = "20151208_143302"
	migration.Register("Tutorpurchasetimer_20151208_143302", m)
}

// Run the migrations
func (m *Tutorpurchasetimer_20151208_143302) Up() {
	m.Sql("ALTER TABLE purchases ADD COLUMN timer_time_start DATETIME")
}

// Reverse the migrations
func (m *Tutorpurchasetimer_20151208_143302) Down() {
	m.Sql("ALTER TABLE purchases DROP COLUMN timer_time_start")
}
