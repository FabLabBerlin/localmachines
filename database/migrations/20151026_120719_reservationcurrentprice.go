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
	m.Sql("ALTER TABLE reservations ADD COLUMN current_price double unsigned")
	m.Sql("ALTER TABLE reservations ADD COLUMN current_price_currency varchar(10)")
	m.Sql("ALTER TABLE reservations ADD COLUMN current_price_unit varchar(100)")
}

// Reverse the migrations
func (m *Reservationcurrentprice_20151026_120719) Down() {
	m.Sql("ALTER TABLE reservations DROP COLUMN current_price")
	m.Sql("ALTER TABLE reservations DROP COLUMN current_price_currency")
	m.Sql("ALTER TABLE reservations DROP COLUMN current_price_unit")
}
