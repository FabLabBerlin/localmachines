package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Purchases_20151119_115310 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Purchases_20151119_115310{}
	m.Created = "20151119_115310"
	migration.Register("Purchases_20151119_115310", m)
}

// Run the migrations
func (m *Purchases_20151119_115310) Up() {
	m.SQL(`CREATE TABLE purchases (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		type varchar(100) NOT NULL,
		product_id int(11) unsigned,
		created datetime,
		user_id int(11) unsigned NOT NULL,
		time_start datetime DEFAULT NULL,
		time_end datetime DEFAULT NULL,
		quantity double NOT NULL,
		price_per_unit double,
		price_unit varchar(100),
		vat double,
		activation_running tinyint(1),
		reservation_disabled tinyint(1),
		machine_id int(11) unsigned,
		PRIMARY KEY (id)
	)`)
	m.SQL(`
		INSERT INTO purchases ( TYPE, product_id, created, user_id, time_start, time_end, quantity, price_per_unit, price_unit, vat, activation_running, reservation_disabled, machine_id )
		SELECT 'activation',
		       NULL,
		       time_start,
		       user_id,
		       time_start,
		       time_end,
		       time_total / 60,
		       current_machine_price,
		       current_machine_price_unit,
		       vat_rate,
		       active,
		       NULL,
		       machine_id
		FROM activations
		WHERE current_machine_price_unit = "minute"
		UNION
		SELECT 'activation',
		       NULL,
		       time_start,
		       user_id,
		       time_start,
		       time_end,
		       time_total / 3600,
		       current_machine_price,
		       current_machine_price_unit,
		       vat_rate,
		       active,
		       NULL,
		       machine_id
		FROM activations
		WHERE current_machine_price_unit = "hour"
	`)
	m.SQL(`
		INSERT INTO purchases ( TYPE, product_id, created, user_id, time_start, time_end, quantity, price_per_unit, price_unit, vat, activation_running, reservation_disabled, machine_id )
		SELECT 'reservation',
		       NULL,
		       created,
		       user_id,
		       time_start,
		       time_end,
		       TIME_TO_SEC(TIMEDIFF(time_end, time_start)) / 1800,
		       current_price / 2,
		       current_price_unit,
		       NULL,
		       NULL,
		       disabled,
		       machine_id
		FROM reservations
	`)
}

// Reverse the migrations
func (m *Purchases_20151119_115310) Down() {
	m.SQL("DROP TABLE purchases")
}
