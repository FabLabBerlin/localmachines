package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Tutorproduct_20151126_120015 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Tutorproduct_20151126_120015{}
	m.Created = "20151126_120015"
	migration.Register("Tutorproduct_20151126_120015", m)
}

// Run the migrations
func (m *Tutorproduct_20151126_120015) Up() {
	m.SQL("ALTER TABLE products ADD COLUMN user_id int(11)")
	m.SQL("ALTER TABLE products ADD COLUMN machine_skills varchar(255)")
}

// Reverse the migrations
func (m *Tutorproduct_20151126_120015) Down() {
	m.SQL("ALTER TABLE products DROP COLUMN machine_skills")
	m.SQL("ALTER TABLE products DROP COLUMN user_id")
}
