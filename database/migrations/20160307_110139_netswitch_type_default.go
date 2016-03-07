package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type NetswitchTypeDefault_20160307_110139 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &NetswitchTypeDefault_20160307_110139{}
	m.Created = "20160307_110139"
	migration.Register("NetswitchTypeDefault_20160307_110139", m)
}

// Run the migrations
func (m *NetswitchTypeDefault_20160307_110139) Up() {
	m.SQL("UPDATE machines SET netswitch_type = 'mfi' WHERE netswitch_xmpp = 1")
}

// Reverse the migrations
func (m *NetswitchTypeDefault_20160307_110139) Down() {
	m.SQL("UPDATE machines SET netswitch_type = NULL WHERE netswitch_xmpp = 1")
}
