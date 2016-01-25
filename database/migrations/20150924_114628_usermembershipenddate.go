package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Usermembershipenddate_20150924_114628 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Usermembershipenddate_20150924_114628{}
	m.Created = "20150924_114628"
	migration.Register("Usermembershipenddate_20150924_114628", m)
}

// Run the migrations
func (m *Usermembershipenddate_20150924_114628) Up() {
	m.SQL(`
UPDATE user_membership
SET end_date = DATE_ADD(start_date, INTERVAL
                          (SELECT duration_months
                           FROM membership
                           WHERE membership.id = user_membership.membership_id) MONTH)
	`)

}

// Reverse the migrations
func (m *Usermembershipenddate_20150924_114628) Down() {
	m.SQL("UPDATE user_membership SET end_date = NULL")
}
