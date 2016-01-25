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
	m.SQL("ALTER TABLE activations ADD COLUMN current_machine_price double unsigned")
	m.SQL("ALTER TABLE activations ADD COLUMN current_machine_price_currency varchar(10)")
	m.SQL("ALTER TABLE activations ADD COLUMN current_machine_price_unit varchar(100)")

	// Fill the new fields of old values
	m.SQL("UPDATE activations a JOIN machines m ON a.machine_id = m.id " +
		"SET a.current_machine_price=m.price, " +
		"a.current_machine_price_currency='€', " +
		"a.current_machine_price_unit=m.price_unit")
}

// Reverse the migrations
func (m *Activationcurrentprice_20151023_094534) Down() {
	m.SQL("ALTER TABLE activations DROP COLUMN current_machine_price")
	m.SQL("ALTER TABLE activations DROP COLUMN current_machine_price_currency")
	m.SQL("ALTER TABLE activations DROP COLUMN current_machine_price_unit")
}
