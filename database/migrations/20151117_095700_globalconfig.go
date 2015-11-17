package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Globalconfig_20151117_095700 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Globalconfig_20151117_095700{}
	m.Created = "20151117_095700"
	migration.Register("Globalconfig_20151117_095700", m)
}

// Run the migrations
func (m *Globalconfig_20151117_095700) Up() {
	m.Sql(`CREATE TABLE global_config (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(100),
		value text,
		PRIMARY KEY (id)
	)`)

}

// Reverse the migrations
func (m *Globalconfig_20151117_095700) Down() {
	m.Sql("DROP TABLE global_config")

}
