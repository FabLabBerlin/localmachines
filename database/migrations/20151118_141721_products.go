package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Products_20151118_141721 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Products_20151118_141721{}
	m.Created = "20151118_141721"
	migration.Register("Products_20151118_141721", m)
}

// Run the migrations
func (m *Products_20151118_141721) Up() {
	m.Sql(`CREATE TABLE products (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		type varchar(100),
		name varchar(100),
		price double unsigned,
		price_unit varchar(100),
		PRIMARY KEY (id)
	)`)
}

// Reverse the migrations
func (m *Products_20151118_141721) Down() {
	m.Sql("DROP TABLE products")
}
