package hexaswitch

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/hexabus"
	"time"
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

// Turn the switch on, the IP address of the switch is being retrieved
// from the database
func On(machineId int) error {

	// Attempt to get switch IP connected to the machine
	switchIp, err := getSwitchIp(machineId)
	if err != nil {
		return err
	}

	// No error, continue using the switch IP got
	beego.Info("Turning ON machine with ID",
		machineId, "and switch IP", switchIp)

	// Attempt to turn the switch on
	err = setSwitchState(true, switchIp)
	if err != nil {
		return err
	}

	// All is fine, clean exit function
	return nil
}

// Turn the switch off, the IP address of the switch is being retrieved
// from the database
func Off(machineId int) error {

	// Attempt to get switch IP connected to the machine
	switchIp, err := getSwitchIp(machineId)
	if err != nil {
		return err
	}

	// No error, use the switch IP
	beego.Info("Turning OFF machine with ID",
		machineId, "and switch IP", switchIp)

	// Attempt to turn the switch off
	err = setSwitchState(false, switchIp)
	if err != nil {
		return err
	}

	// All is fine, return no-error
	return nil
}

func setSwitchState(switchState bool, switchIp string) error {

	// Create write packet to switch on and off
	var writePacket hexabus.WritePacket = hexabus.WritePacket{hexabus.FLAG_NONE,
		1, hexabus.DTYPE_BOOL, switchState}

	switchStateStr := "On"
	if !switchState {
		switchStateStr = "Off"
	}
	beego.Info("Sending hexaswitch packet", switchStateStr)

	// Register time before sending the packet
	timeBeforeSendingPacket := time.Now()

	// Send packet to switch
	err := writePacket.Send(switchIp)

	// Register time after sending the packet
	timeResponseReceived := time.Now()
	beego.Info("Communicating with the switch took",
		timeResponseReceived.Sub(timeBeforeSendingPacket).Seconds(), "seconds")

	// Check for errors
	if err != nil {
		return err
	}

	// Attempt to read the switch in order to get current state
	var queryPacket hexabus.QueryPacket = hexabus.QueryPacket{hexabus.FLAG_NONE, 1}
	var bytes []byte
	bytes, err = queryPacket.Send(switchIp)
	if err != nil {
		beego.Error(err)
	}
	beego.Trace("Query bytes", bytes)
	// Query bytes [72 88 48 67 1 0 0 0 0 1 1 1 192 215] // switching on
	// Query bytes [72 88 48 67 1 0 0 0 0 1 1 0 209 94] // switching off
	ptype, err := hexabus.PacketType(bytes)
	if err != nil {
		return err
	}
	beego.Trace("Packet type", ptype)

	if ptype == hexabus.PTYPE_INFO {
		beego.Info("Received hexabus info packet")
		infoPacket := hexabus.InfoPacket{}
		err = infoPacket.Decode(bytes)
		if err != nil {
			beego.Error("Failed to decode hexabus info packet", err)
		}

		// Expecting boolean value
		switchState := infoPacket.Data
		beego.Trace("Info pack switch state:", switchState)
	}

	// No errors so far, return no-error
	return nil
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
