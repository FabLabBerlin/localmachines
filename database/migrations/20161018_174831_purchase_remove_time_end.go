package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PurchaseRemoveTimeEnd_20161018_174831 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PurchaseRemoveTimeEnd_20161018_174831{}
	m.Created = "20161018_174831"
	migration.Register("PurchaseRemoveTimeEnd_20161018_174831", m)
}

// Run the migrations
func (m *PurchaseRemoveTimeEnd_20161018_174831) Up() {
	m.SQL("ALTER TABLE purchases DROP COLUMN time_end")
}

// Reverse the migrations
func (m *PurchaseRemoveTimeEnd_20161018_174831) Down() {
	m.SQL("ALTER TABLE purchases ADD COLUMN time_end datetime AFTER time_start")
	m.SQL("UPDATE purchases SET time_end = date_add(time_start, interval quantity minute) WHERE price_unit = 'minute' AND running = 0")
	m.SQL("UPDATE purchases SET time_end = date_add(time_start, interval (quantity / 2) hour) WHERE price_unit = '30 minutes' AND running = 0")
	m.SQL("UPDATE purchases SET time_end = date_add(time_start, interval quantity hour) WHERE price_unit = 'hour' AND running = 0")
	m.SQL("UPDATE purchases SET time_end = date_add(time_start, interval quantity month) WHERE price_unit = 'month' AND running = 0")
}
