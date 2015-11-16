package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Datalog_20151113_123217 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Datalog_20151113_123217{}
	m.Created = "20151113_123217"
	migration.Register("Datalog_20151113_123217", m)
}

// Run the migrations
func (m *Datalog_20151113_123217) Up() {
	m.Sql(`
		CREATE TABLE data_log (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			changed_table varchar(100) NOT NULL DEFAULT '',
			before_json text,
			after_json text,
			created datetime NOT NULL,
			hash text NOT NULL,
			PRIMARY KEY (id)
		)
	`)
}

// Reverse the migrations
func (m *Datalog_20151113_123217) Down() {
	m.Sql("DROP TABLE data_log")
}
