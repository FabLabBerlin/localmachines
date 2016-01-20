package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MergeHostsAndLocations_20160120_163019 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MergeHostsAndLocations_20160120_163019{}
	m.Created = "20160120_163019"
	migration.Register("MergeHostsAndLocations_20160120_163019", m)
}

// Run the migrations
func (m *MergeHostsAndLocations_20160120_163019) Up() {
	m.Sql("DROP TABLE hosts")
	m.Sql("ALTER TABLE locations CHANGE name title varchar(100)")
	m.Sql("ALTER TABLE locations ADD COLUMN first_name varchar(100)")
	m.Sql("ALTER TABLE locations ADD COLUMN last_name varchar(100)")
	m.Sql("ALTER TABLE locations ADD COLUMN email varchar(100)")
	m.Sql("ALTER TABLE locations ADD COLUMN city varchar(100)")
	m.Sql("ALTER TABLE locations ADD COLUMN organization varchar(100)")
	m.Sql("ALTER TABLE locations ADD COLUMN phone varchar(100)")
	m.Sql("ALTER TABLE locations ADD COLUMN comments varchar(100)")
	m.Sql("ALTER TABLE locations ADD COLUMN approved tinyint(1)")
	m.Sql("UPDATE locations SET approved = 1 WHERE id = 1")
}

// Reverse the migrations
func (m *MergeHostsAndLocations_20160120_163019) Down() {
	m.Sql(`
        CREATE TABLE hosts (
            id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
            first_name varchar(100) NOT NULL,
            last_name varchar(100) NOT NULL,
            email varchar(100) NOT NULL,
            location varchar(100) NOT NULL,
            organization varchar(100) NOT NULL,
            phone varchar(100) NOT NULL,
            comments text NOT NULL,
            PRIMARY KEY (id)
    )`)
	m.Sql("ALTER TABLE locations DROP COLUMN first_name")
	m.Sql("ALTER TABLE locations DROP COLUMN last_name")
	m.Sql("ALTER TABLE locations DROP COLUMN email")
	m.Sql("ALTER TABLE locations DROP COLUMN city")
	m.Sql("ALTER TABLE locations DROP COLUMN organization")
	m.Sql("ALTER TABLE locations DROP COLUMN phone")
	m.Sql("ALTER TABLE locations DROP COLUMN comments")
	m.Sql("ALTER TABLE locations DROP COLUMN approved")
	m.Sql("ALTER TABLE locations CHANGE title name varchar(100)")
}
