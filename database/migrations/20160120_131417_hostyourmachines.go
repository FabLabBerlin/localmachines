package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Hostyourmachines_20160120_131417 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Hostyourmachines_20160120_131417{}
	m.Created = "20160120_131417"
	migration.Register("Hostyourmachines_20160120_131417", m)
}

// Run the migrations
func (m *Hostyourmachines_20160120_131417) Up() {
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
}

// Reverse the migrations
func (m *Hostyourmachines_20160120_131417) Down() {
	m.Sql("DROP TABLE hosts")
}
