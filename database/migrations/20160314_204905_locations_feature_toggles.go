package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type LocationsFeatureToggles_20160314_204905 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &LocationsFeatureToggles_20160314_204905{}
	m.Created = "20160314_204905"
	migration.Register("LocationsFeatureToggles_20160314_204905", m)
}

// Run the migrations
func (m *LocationsFeatureToggles_20160314_204905) Up() {
	m.SQL("ALTER TABLE locations ADD COLUMN feature_coworking TINYINT(1)")
	m.SQL("ALTER TABLE locations ADD COLUMN feature_setup_time TINYINT(1)")
	m.SQL("ALTER TABLE locations ADD COLUMN feature_spaces TINYINT(1)")
	m.SQL("ALTER TABLE locations ADD COLUMN feature_tutoring TINYINT(1)")
	m.SQL("UPDATE locations SET feature_coworking = 1 WHERE id = 1")
	m.SQL("UPDATE locations SET feature_setup_time = 1 WHERE id = 1")
	m.SQL("UPDATE locations SET feature_spaces = 1 WHERE id = 1")
	m.SQL("UPDATE locations SET feature_tutoring = 1 WHERE id = 1")
}

// Reverse the migrations
func (m *LocationsFeatureToggles_20160314_204905) Down() {
	m.SQL("ALTER TABLE locations DROP COLUMN feature_coworking")
	m.SQL("ALTER TABLE locations DROP COLUMN feature_setup_time")
	m.SQL("ALTER TABLE locations DROP COLUMN feature_spaces")
	m.SQL("ALTER TABLE locations DROP COLUMN feature_tutoring")
}
