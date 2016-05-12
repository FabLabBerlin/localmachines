package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MonthlyEarningsOneMonth_20160512_154039 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MonthlyEarningsOneMonth_20160512_154039{}
	m.Created = "20160512_154039"
	migration.Register("MonthlyEarningsOneMonth_20160512_154039", m)
}

// Run the migrations
func (m *MonthlyEarningsOneMonth_20160512_154039) Up() {
	m.SQL("ALTER TABLE monthly_earnings CHANGE month_from month TINYINT UNSIGNED")
	m.SQL("ALTER TABLE monthly_earnings CHANGE year_from year SMALLINT UNSIGNED")
	m.SQL("ALTER TABLE monthly_earnings DROP COLUMN month_to")
	m.SQL("ALTER TABLE monthly_earnings DROP COLUMN year_to")
}

// Reverse the migrations
func (m *MonthlyEarningsOneMonth_20160512_154039) Down() {
	m.SQL("ALTER TABLE monthly_earnings CHANGE month month_from TINYINT UNSIGNED")
	m.SQL("ALTER TABLE monthly_earnings CHANGE year year_from SMALLINT UNSIGNED")
	m.SQL("ALTER TABLE monthly_earnings ADD COLUMN month_to TINYINT UNSIGNED AFTER year_from")
	m.SQL("ALTER TABLE monthly_earnings ADD COLUMN year_to SMALLINT UNSIGNED AFTER month_to")
	m.SQL("UPDATE monthly_earnings SET month_to = month_from")
	m.SQL("UPDATE monthly_earnings SET year_to = year_from")
}
