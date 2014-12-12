package hexaswitch

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/morriswinkler/hexabus"
)

var switchAddress string = "[fafa::50:c4ff:fe04:8390]"
var tableName string = "hexaswitch"

func init() {

	// We launch Install every time on init now.
	// To be removed later
	orm.RegisterModel(new(HexaSwitch))
	//Install()
}

// Hexaswitch plugin DB table model
type HexaSwitch struct {
	Id        int `orm:"auto"`
	MachineId int
	SwitchIp  string `orm:"size(100)"`
}

// Creates a MySQL table for machine and switch IP mappings
func Install() {
	o := orm.NewOrm()
	res, err := o.Raw("CREATE TABLE IF NOT EXISTS hexaswitch (id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY, machine_id INT NOT NULL, switch_ip VARCHAR(100) NOT NULL)").Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		beego.Trace("MySQL rows affected:", num)
	} else {
		beego.Error(err)
	}
}

func On(machineId int) {

	switchIp, err := getSwitchIp(machineId)
	if err == nil {
		beego.Trace("Turning ON machine with ID", machineId, "switch IP", switchIp)
		go setSwitchState(true, switchIp)
	} else {
		beego.Error(err)
	}
}

func Off(machineId int) {
	switchIp, err := getSwitchIp(machineId)

	if err == nil {
		beego.Trace("Turning OFF machine with ID", machineId, "switch IP", switchIp)
		go setSwitchState(false, switchIp)
	} else {
		beego.Error(err)
	}
}

func setSwitchState(switchState bool, switchIp string) {

	// Create write packet to switch on and off
	var wPack hexabus.WritePacket = hexabus.WritePacket{hexabus.FLAG_NONE,
		1, hexabus.DTYPE_BOOL, switchState}

	switchStateStr := "Off"
	if !switchState {
		switchStateStr = "On"
	}
	beego.Trace("Hexaswitch", switchStateStr)

	// Send packet to switch
	err := wPack.Send(switchIp)
	if err != nil {
		fmt.Println(err)
	}
}

func getSwitchIp(machineId int) (string, error) {
	switchModel := HexaSwitch{}
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM hexaswitch WHERE machine_id = ? LIMIT 0, 1", machineId).QueryRow(&switchModel)
	if err != nil {
		return "", err
	}
	return switchModel.SwitchIp, nil
}
