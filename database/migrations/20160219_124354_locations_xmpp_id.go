package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type LocationsXmppId_20160219_124354 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &LocationsXmppId_20160219_124354{}
	m.Created = "20160219_124354"
	migration.Register("LocationsXmppId_20160219_124354", m)
}

// Run the migrations
func (m *LocationsXmppId_20160219_124354) Up() {
	m.SQL("ALTER TABLE locations ADD COLUMN xmpp_id VARCHAR(255)")
}

// Reverse the migrations
func (m *LocationsXmppId_20160219_124354) Down() {
	m.SQL("ALTER TABLE locations DROP COLUMN xmpp_id")
}
