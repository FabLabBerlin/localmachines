package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Coupons_20160421_144423 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Coupons_20160421_144423{}
	m.Created = "20160421_144423"
	migration.Register("Coupons_20160421_144423", m)
}

// Run the migrations
func (m *Coupons_20160421_144423) Up() {
	m.SQL(`
		CREATE TABLE coupons (
			id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
			location_id INT(11) UNSIGNED,
			code VARCHAR(100),
			user_id INT(11) UNSIGNED,
			value DOUBLE,
			PRIMARY KEY (id)
	)`)
	m.SQL(`
		CREATE TABLE coupon_usages (
			id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
			coupon_id INT(11) UNSIGNED,
			value DOUBLE,
			month TINYINT,
			year SMALLINT,
			PRIMARY KEY (id)
	)`)
}

// Reverse the migrations
func (m *Coupons_20160421_144423) Down() {
	m.SQL("DROP TABLE coupons")
	m.SQL("DROP TABLE coupon_usages")
}
