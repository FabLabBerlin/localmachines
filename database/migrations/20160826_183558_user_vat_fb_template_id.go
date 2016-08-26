package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserVatFbTemplateId_20160826_183558 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserVatFbTemplateId_20160826_183558{}
	m.Created = "20160826_183558"
	migration.Register("UserVatFbTemplateId_20160826_183558", m)
}

// Run the migrations
func (m *UserVatFbTemplateId_20160826_183558) Up() {
	m.SQL("ALTER TABLE user ADD COLUMN fastbill_template_id int(11) unsigned")
	m.SQL("ALTER TABLE user ADD COLUMN eu_delivery tinyint(1)")
}

// Reverse the migrations
func (m *UserVatFbTemplateId_20160826_183558) Down() {
	m.SQL("ALTER TABLE user DROP COLUMN fastbill_template_id")
	m.SQL("ALTER TABLE user DROP COLUMN eu_delivery")

}
