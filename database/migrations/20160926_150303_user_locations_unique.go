package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserLocationsUnique_20160926_150303 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserLocationsUnique_20160926_150303{}
	m.Created = "20160926_150303"
	migration.Register("UserLocationsUnique_20160926_150303", m)
}

// Run the migrations
func (m *UserLocationsUnique_20160926_150303) Up() {
	m.SQL("ALTER TABLE user_locations ADD UNIQUE unique_user_locations (user_id, location_id)")
}

// Reverse the migrations
func (m *UserLocationsUnique_20160926_150303) Down() {
	m.SQL("DROP INDEX unique_user_locations ON user_locations")
}
