package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MergeMachinesAndNetswitches_20160219_133411 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MergeMachinesAndNetswitches_20160219_133411{}
	m.Created = "20160219_133411"
	migration.Register("MergeMachinesAndNetswitches_20160219_133411", m)
}

// Run the migrations
func (m *MergeMachinesAndNetswitches_20160219_133411) Up() {
	m.SQL(`
		ALTER TABLE machines 
		  ADD COLUMN netswitch_url_on VARCHAR(255),
		  ADD COLUMN netswitch_url_off VARCHAR(255),
		  ADD COLUMN netswitch_host VARCHAR(255),
		  ADD COLUMN netswitch_sensor_port INT(5),
		  ADD COLUMN netswitch_xmpp TINYINT(1)
	`)
	m.SQL(`
		UPDATE machines
		       JOIN netswitch
		         ON machines.id = netswitch.machine_id
		SET    netswitch_url_on = netswitch.url_on,
		       netswitch_url_off = netswitch.url_off,
		       netswitch_host = netswitch.host,
		       netswitch_sensor_port = netswitch.sensor_port,
		       netswitch_xmpp = netswitch.xmpp
	`)
	m.SQL(`DROP TABLE netswitch`)
}

// Reverse the migrations
func (m *MergeMachinesAndNetswitches_20160219_133411) Down() {
	m.SQL(`
		CREATE TABLE netswitch 
		  ( 
		     id          INT(11) UNSIGNED NOT NULL auto_increment,
		     location_id INT(11) UNSIGNED NOT NULL,
		     machine_id  INT(11) UNSIGNED NOT NULL,
		     url_on      VARCHAR(255),
		     url_off     VARCHAR(255),
		     host        VARCHAR(255),
		     sensor_port INT(5),
		     xmpp        TINYINT(1),
		     PRIMARY KEY (id)
		  )
	`)
	m.SQL(`
		UPDATE netswitch
		       JOIN machines
		         ON machines.id = netswitch.machine_id
		SET    netswitch.url_on = netswitch_url_on,
		       netswitch.url_off = netswitch_url_off,
		       netswitch.host = netswitch_host,
		       netswitch.sensor_port = netswitch_sensor_port,
		       netswitch.xmpp = netswitch_xmpp
	`)
	m.SQL(`
		ALTER TABLE machines
		  DROP COLUMN netswitch_url_on,
		  DROP COLUMN netswitch_url_off,
		  DROP COLUMN netswitch_host,
		  DROP COLUMN netswitch_sensor_port,
		  DROP COLUMN netswitch_xmpp
	`)
}
