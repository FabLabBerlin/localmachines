package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MonthlyEarnings_20160304_113146 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MonthlyEarnings_20160304_113146{}
	m.Created = "20160304_113146"
	migration.Register("MonthlyEarnings_20160304_113146", m)
}

// Run the migrations
func (m *MonthlyEarnings_20160304_113146) Up() {
	m.SQL("RENAME TABLE invoices TO monthly_earnings")
	m.SQL("ALTER TABLE monthly_earnings ADD COLUMN month_from TINYINT UNSIGNED AFTER id")
	m.SQL("ALTER TABLE monthly_earnings ADD COLUMN year_from SMALLINT UNSIGNED AFTER month_from")
	m.SQL("ALTER TABLE monthly_earnings ADD COLUMN month_to TINYINT UNSIGNED AFTER year_from")
	m.SQL("ALTER TABLE monthly_earnings ADD COLUMN year_to SMALLINT UNSIGNED AFTER month_to")
	m.SQL("UPDATE monthly_earnings SET month_from = MONTH(period_from)")
	m.SQL("UPDATE monthly_earnings SET year_from = YEAR(period_from)")
	m.SQL("UPDATE monthly_earnings SET month_to = MONTH(period_to)")
	m.SQL("UPDATE monthly_earnings SET year_to = YEAR(period_to)")
}

// Reverse the migrations
func (m *MonthlyEarnings_20160304_113146) Down() {
	m.SQL("ALTER TABLE monthly_earnings DROP COLUMN month_from")
	m.SQL("ALTER TABLE monthly_earnings DROP COLUMN year_from")
	m.SQL("ALTER TABLE monthly_earnings DROP COLUMN month_to")
	m.SQL("ALTER TABLE monthly_earnings DROP COLUMN year_to")
	m.SQL("RENAME TABLE monthly_earnings TO invoices")
}
