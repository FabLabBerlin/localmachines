package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Activationfeedback_20150908_145935 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Activationfeedback_20150908_145935{}
	m.Created = "20150908_145935"
	migration.Register("Activationfeedback_20150908_145935", m)
}

// Run the migrations
func (m *Activationfeedback_20150908_145935) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`
		CREATE TABLE activation_feedback (
			id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
			activation_id int(11) NOT NULL,
			satisfaction varchar(100) DEFAULT NULL,
			PRIMARY KEY (id)
	)`)
}

// Reverse the migrations
func (m *Activationfeedback_20150908_145935) Down() {
	m.SQL("DROP TABLE activation_feedback")
}
