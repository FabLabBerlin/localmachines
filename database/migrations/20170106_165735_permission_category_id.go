package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type PermissionCategoryId_20170106_165735 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &PermissionCategoryId_20170106_165735{}
	m.Created = "20170106_165735"
	migration.Register("PermissionCategoryId_20170106_165735", m)
}

// Run the migrations
func (m *PermissionCategoryId_20170106_165735) Up() {
	m.SQL(`
		CREATE TABLE permission_new (
			id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
			location_id INT(11) UNSIGNED,
			user_id INT(11) UNSIGNED,
			category_id INT(11) UNSIGNED,
			PRIMARY KEY (id)
	)`)
	m.SQL(`
INSERT INTO permission_new
SELECT permission.id,
       location_id,
       permission.user_id,
       type_id
FROM   permission
       JOIN machines
         ON machines.id = permission.machine_id
GROUP  BY permission.user_id,
          type_id 
	`)
	m.SQL(`RENAME TABLE permission TO permission_old`)
	m.SQL(`RENAME TABLE permission_new TO permission`)
}

// Reverse the migrations
func (m *PermissionCategoryId_20170106_165735) Down() {
	m.SQL(`DROP TABLE permission`)
	m.SQL(`RENAME TABLE permission_old TO permission`)
}
