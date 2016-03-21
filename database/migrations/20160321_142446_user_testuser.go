package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserTestuser_20160321_142446 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserTestuser_20160321_142446{}
	m.Created = "20160321_142446"
	migration.Register("UserTestuser_20160321_142446", m)
}

// Run the migrations
func (m *UserTestuser_20160321_142446) Up() {
	m.SQL("ALTER TABLE user ADD COLUMN test_user TINYINT(1)")
}

// Reverse the migrations
func (m *UserTestuser_20160321_142446) Down() {
	m.SQL("ALTER TABLE user DROP COLUMN test_user")
}
