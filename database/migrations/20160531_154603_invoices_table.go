package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type InvoicesTable_20160531_154603 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &InvoicesTable_20160531_154603{}
	m.Created = "20160531_154603"
	migration.Register("InvoicesTable_20160531_154603", m)
}

// Run the migrations
func (m *InvoicesTable_20160531_154603) Up() {
	m.SQL(`
		CREATE TABLE invoices (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			location_id int(11) unsigned,
			fastbill_id int(11) unsigned,
			fastbill_no varchar(100),
			month tinyint unsigned,
			year smallint unsigned,
			customer_id int(11) unsigned,
			customer_no int(11) unsigned,
			user_id int(11) unsigned,
			status varchar(20),
			canceled tinyint(1) DEFAULT 0,
			PRIMARY KEY (id)
	)`)
}

// Reverse the migrations
func (m *InvoicesTable_20160531_154603) Down() {
	m.SQL("DROP TABLE invoices")
}
