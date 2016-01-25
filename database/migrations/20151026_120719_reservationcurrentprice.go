package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Reservationcurrentprice_20151026_120719 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Reservationcurrentprice_20151026_120719{}
	m.Created = "20151026_120719"
	migration.Register("Reservationcurrentprice_20151026_120719", m)
}

// Run the migrations
func (m *Reservationcurrentprice_20151026_120719) Up() {
	m.SQL("ALTER TABLE reservations ADD COLUMN current_price double unsigned")
	m.SQL("ALTER TABLE reservations ADD COLUMN current_price_currency varchar(10)")
	m.SQL("ALTER TABLE reservations ADD COLUMN current_price_unit varchar(100)")

	// Fill the new fields of old records
	m.SQL("UPDATE reservations r JOIN machines m ON r.machine_id = m.id " +
		"SET r.current_price=m.reservation_price_hourly, " +
		"r.current_price_currency='€', " +
		"r.current_price_unit='30 minutes'")
}

// Reverse the migrations
func (m *Reservationcurrentprice_20151026_120719) Down() {
	m.SQL("ALTER TABLE reservations DROP COLUMN current_price")
	m.SQL("ALTER TABLE reservations DROP COLUMN current_price_currency")
	m.SQL("ALTER TABLE reservations DROP COLUMN current_price_unit")
}
