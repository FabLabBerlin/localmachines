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
	m.SQL(`CREATE TABLE settings (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(100) NOT NULL,
		value_int int(11),
		value_string text,
		value_float double,
		PRIMARY KEY (id)
	)`)

}

// Reverse the migrations
func (m *Globalconfig_20151117_095700) Down() {
	m.SQL("DROP TABLE settings")
}
