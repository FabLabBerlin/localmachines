package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Activationcurrentprice_20151023_094534 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Activationcurrentprice_20151023_094534{}
	m.Created = "20151023_094534"
	migration.Register("Activationcurrentprice_20151023_094534", m)
}

// Run the migrations
func (m *Activationcurrentprice_20151023_094534) Up() {
	m.Sql("ALTER TABLE activations ADD COLUMN current_machine_price double unsigned")
	m.Sql("ALTER TABLE activations ADD COLUMN current_machine_price_currency varchar(10)")
	m.Sql("ALTER TABLE activations ADD COLUMN current_machine_price_unit varchar(100)")
}

// Reverse the migrations
func (m *Activationcurrentprice_20151023_094534) Down() {
	m.Sql("ALTER TABLE activations DROP COLUMN current_machine_price")
	m.Sql("ALTER TABLE activations DROP COLUMN current_machine_price_currency")
	m.Sql("ALTER TABLE activations DROP COLUMN current_machine_price_unit")
}
