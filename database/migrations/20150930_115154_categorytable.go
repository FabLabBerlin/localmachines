package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Categorytable_20150930_115154 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Categorytable_20150930_115154{}
	m.Created = "20150930_115154"
	migration.Register("Categorytable_20150930_115154", m)
}

// Run the migrations
func (m *Categorytable_20150930_115154) Up() {
	m.Sql(`
		CREATE TABLE categories (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			name varchar(100) DEFAULT NULL,
			shortname varchar(100),
			description text,
			image varchar(255), 
			PRIMARY KEY (id)
	)`)
	m.Sql(`ALTER TABLE machines DROP COLUMN category`)
	m.Sql(`ALTER TABLE machines ADD COLUMN category_id int(11)`)
}

// Reverse the migrations
func (m *Categorytable_20150930_115154) Down() {
	m.Sql("DROP TABLE categories")
	m.Sql(`ALTER TABLE machines ADD COLUMN category varchar(30)`)
	m.Sql(`ALTER TABLE machines DROP COLUMN category_id`)
}
