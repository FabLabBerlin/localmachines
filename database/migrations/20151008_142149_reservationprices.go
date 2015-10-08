package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Reservationprices_20151008_142149 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Reservationprices_20151008_142149{}
	m.Created = "20151008_142149"
	migration.Register("Reservationprices_20151008_142149", m)
}

// Run the migrations
func (m *Reservationprices_20151008_142149) Up() {
	m.Sql("ALTER TABLE machines ADD COLUMN reservation_price_start double unsigned")
	m.Sql("ALTER TABLE machines ADD COLUMN reservation_price_hourly double unsigned")
}

// Reverse the migrations
func (m *Reservationprices_20151008_142149) Down() {
	m.Sql("ALTER TABLE machines DROP COLUMN reservation_price_start")
	m.Sql("ALTER TABLE machines DROP COLUMN reservation_price_hourly")
}
