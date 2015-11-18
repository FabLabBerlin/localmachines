package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Purchases_20151118_162139 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Purchases_20151118_162139{}
	m.Created = "20151118_162139"
	migration.Register("Purchases_20151118_162139", m)
}

// Run the migrations
func (m *Purchases_20151118_162139) Up() {
	m.Sql("ALTER TABLE activations ADD COLUMN purchase_id int(11) unsigned")
	m.Sql("ALTER TABLE machines ADD COLUMN product_id int(11) unsigned")
	m.Sql("ALTER TABLE machines ADD COLUMN reservation_product_id int(11) unsigned")
	m.Sql("UPDATE machines SET product_id = id")
	m.Sql("UPDATE machines SET reservation_product_id = id + 100")
	m.Sql(`CREATE TABLE purchases (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		product_id int(11) unsigned NOT NULL,
		user_id int(11) unsigned NOT NULL,
		time_start datetime DEFAULT NULL,
		time_end datetime DEFAULT NULL,
		quantity double NOT NULL,
		price double,
		price_unit varchar(100),
		PRIMARY KEY (id)
	)`)
	m.Sql(`
		INSERT INTO products (id, type, name, price, price_unit)
		SELECT id,
		       "machine",
		       name,
		       price,
		       price_unit
		FROM machines
		UNION
		SELECT id + 100,
		       "machine_reservation",
		       "Reservation (" + name + ")",
		       price,
		       price_unit
		FROM machines`)
	m.Sql(`
		INSERT INTO purchases ( id, product_id, user_id, time_start, time_end, quantity, price, price_unit )
		SELECT id,
		       machine_id,
		       user_id,
		       time_start,
		       time_end,
		       time_total,
		       current_machine_price,
		       current_machine_price_unit
		FROM activations
	`)
	m.Sql("UPDATE activations SET purchase_id = id")
	m.Sql(`
		INSERT INTO purchases ( id, product_id, user_id, time_start, time_end, quantity, price, price_unit )
		SELECT id + 10000,
		       machine_id + 100,
		       user_id,
		       time_start,
		       time_end,
		       TIME_TO_SEC(TIMEDIFF(time_end, time_start)) / 1800,
		       current_price,
		       current_price_unit
		FROM reservations
	`)
	m.Sql("ALTER TABLE activations MODIFY user_id int(11)")
	m.Sql("ALTER TABLE activations MODIFY time_total int(11)")
}

// Reverse the migrations
func (m *Purchases_20151118_162139) Down() {
	m.Sql("DELETE FROM products")
	m.Sql("DROP TABLE purchases")
	m.Sql("ALTER TABLE activations DROP COLUMN purchase_id")
	m.Sql("ALTER TABLE machines DROP COLUMN product_id")
	m.Sql("ALTER TABLE machines DROP COLUMN reservation_product_id")
}
