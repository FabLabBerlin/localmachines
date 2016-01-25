package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PricePerMonth_20150901_110810 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PricePerMonth_20150901_110810{}
	m.Created = "20150901_110810"
	migration.Register("PricePerMonth_20150901_110810", m)
}

// Run the migrations
func (m *PricePerMonth_20150901_110810) Up() {
	m.SQL("ALTER TABLE membership ADD monthly_price double unsigned NOT NULL AFTER price")
	m.SQL("UPDATE membership SET monthly_price = price WHERE duration = 30 AND unit = 'days'")
	m.SQL("UPDATE membership SET monthly_price = price / 3 WHERE duration = 90 AND unit = 'days'")
	m.SQL("UPDATE membership SET monthly_price = price / 12 WHERE duration = 365 AND unit = 'days'")
	m.SQL("UPDATE membership SET monthly_price = price / duration * 30 WHERE duration <> 30 AND duration <> 90 AND duration <> 365 AND unit = 'days'")
	m.SQL("UPDATE membership SET monthly_price = price WHERE unit <> 'days'")
	m.SQL("ALTER TABLE membership DROP COLUMN price")
}

// Reverse the migrations
func (m *PricePerMonth_20150901_110810) Down() {
	m.SQL("ALTER TABLE membership ADD price double unsigned NOT NULL AFTER monthly_price")
	m.SQL("UPDATE membership SET price = monthly_price WHERE duration = 30 AND unit = 'days'")
	m.SQL("UPDATE membership SET price = monthly_price * 3 WHERE duration = 90 AND unit = 'days'")
	m.SQL("UPDATE membership SET price = monthly_price * 12 WHERE duration = 365 AND unit = 'days'")
	m.SQL("UPDATE membership SET price = monthly_price * duration / 30 WHERE duration <> 30 AND duration <> 90 AND duration <> 365 AND unit = 'days'")
	m.SQL("UPDATE membership SET price = monthly_price  WHERE unit <> 'days'")
	m.SQL("ALTER TABLE membership DROP COLUMN monthly_price")
}
