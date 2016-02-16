package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserLocationsMemberRoles_20160216_131116 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserLocationsMemberRoles_20160216_131116{}
	m.Created = "20160216_131116"
	migration.Register("UserLocationsMemberRoles_20160216_131116", m)
}

// Run the migrations
func (m *UserLocationsMemberRoles_20160216_131116) Up() {
	// All existing non-admins in the system become explicitly members of
	// FabLab Berlin.
	m.SQL(`UPDATE user SET user_role = "member" WHERE user_role <> "admin"`)
	m.SQL(`
		INSERT INTO user_locations
		            (location_id,
		             user_id,
		             user_role,
		             archived)
		SELECT 1,
		       id,
		       "member",
		       0
		FROM   user
		WHERE  user_role <> "admin"
	`)
}

// Reverse the migrations
func (m *UserLocationsMemberRoles_20160216_131116) Down() {
	m.SQL(`DELETE FROM user_locations WHERE user_role <> "admin"`)
}
