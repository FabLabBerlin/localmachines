package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserLocationsAdminRoles_20160210_121712 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserLocationsAdminRoles_20160210_121712{}
	m.Created = "20160210_121712"
	migration.Register("UserLocationsAdminRoles_20160210_121712", m)
}

// Run the migrations
func (m *UserLocationsAdminRoles_20160210_121712) Up() {
	// All existing admins in the system become explicitly admins of
	// FabLab Berlin.
	m.SQL(`
		INSERT INTO user_locations
		            (location_id,
		             user_id,
		             user_role,
		             archived)
		SELECT 1,
		       id,
		       "admin",
		       0
		FROM   user
		WHERE  user_role = "admin"
	`)
}

// Reverse the migrations
func (m *UserLocationsAdminRoles_20160210_121712) Down() {
	m.SQL("DELETE FROM user_locations")
}
