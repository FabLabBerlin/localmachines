package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type LocationsFeatureToggleCoupons_20160422_150028 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &LocationsFeatureToggleCoupons_20160422_150028{}
	m.Created = "20160422_150028"
	migration.Register("LocationsFeatureToggleCoupons_20160422_150028", m)
}

// Run the migrations
func (m *LocationsFeatureToggleCoupons_20160422_150028) Up() {
	m.SQL("ALTER TABLE locations ADD COLUMN feature_coupons TINYINT(1)")
	m.SQL("UPDATE locations SET feature_coupons = 1 WHERE id = 1")
}

// Reverse the migrations
func (m *LocationsFeatureToggleCoupons_20160422_150028) Down() {
	m.SQL("ALTER TABLE locations DROP COLUMN feature_coupons")
}
