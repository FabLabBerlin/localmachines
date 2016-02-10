package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ProductLocationId_20160210_145725 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ProductLocationId_20160210_145725{}
	m.Created = "20160210_145725"
	migration.Register("ProductLocationId_20160210_145725", m)
}

// Run the migrations
func (m *ProductLocationId_20160210_145725) Up() {
	m.SQL(`
	    ALTER TABLE products 
	      ADD COLUMN location_id INT(11) after id
	`)
	m.SQL(`UPDATE products SET location_id = 1`)
}

// Reverse the migrations
func (m *ProductLocationId_20160210_145725) Down() {
	m.SQL("ALTER TABLE products DROP COLUMN location_id")
}
