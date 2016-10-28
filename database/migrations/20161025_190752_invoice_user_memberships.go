package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type InvoiceUserMemberships_20161025_190752 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &InvoiceUserMemberships_20161025_190752{}
	m.Created = "20161025_190752"
	migration.Register("InvoiceUserMemberships_20161025_190752", m)
}

// Run the migrations
func (m *InvoiceUserMemberships_20161025_190752) Up() {
	m.SQL(`
CREATE TABLE user_memberships (
	id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	location_id INT(11) UNSIGNED,
	user_id INT(11) UNSIGNED,
	membership_id INT(11) UNSIGNED,
	start_date DATETIME,
	termination_date DATETIME,
	initial_duration_months INT(11),
	auto_extend TINYINT(1),
	created DATETIME,
	updated DATETIME,
	PRIMARY KEY (id)
)
	`)
	m.SQL(`
CREATE TABLE invoice_user_memberships (
	id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
	location_id INT(11) UNSIGNED,
	user_id INT(11) UNSIGNED,
	membership_id INT(11) UNSIGNED,
	user_membership_id INT(11) UNSIGNED,
	start_date DATETIME,
	termination_date DATETIME,
	initial_duration_months INT(11),
	created DATETIME,
	updated DATETIME,
	invoice_id INT(11),
	invoice_status VARCHAR(100),
	PRIMARY KEY (id)
)
	`)
}

// Reverse the migrations
func (m *InvoiceUserMemberships_20161025_190752) Down() {
	m.SQL("DROP TABLE user_memberships")
	m.SQL("DROP TABLE invoice_user_memberships")
}
