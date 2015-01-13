package hexaswitch

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/hexabus"
	"time"
)

const (
	SWITCH_STATE_ON  = true
	SWITCH_STATE_OFF = false
)

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
func SwitchOn(machineId int) error {

	// Attempt to get switch IP connected to the machine
	switchIp, err := getSwitchIp(machineId)
	if err != nil {
		return err
	}

	// No error, continue using the switch IP got
	beego.Info("Turning ON machine with ID",
		machineId, "and switch IP", switchIp)

	// Attempt to turn the switch on
	err = setSwitchState(SWITCH_STATE_ON, switchIp)
	if err != nil {
		return err
	}

	// All is fine, clean exit function
	return nil
}

// Turn the switch off, the IP address of the switch is being retrieved
// from the database
func SwitchOff(machineId int) error {

	// Attempt to get switch IP connected to the machine
	switchIp, err := getSwitchIp(machineId)
	if err != nil {
		return err
	}

	// No error, use the switch IP
	beego.Info("Turning OFF machine with ID",
		machineId, "and switch IP", switchIp)

	// Attempt to turn the switch off
	err = setSwitchState(SWITCH_STATE_OFF, switchIp)
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

	timeBeforeQuery := time.Now()

	// Short timeout before checking current switch status
	time.Sleep(time.Millisecond * 100)

	// Attempt to read the switch in order to get current state,
	// to check if the write actually worked
	var queryPacket hexabus.QueryPacket = hexabus.QueryPacket{hexabus.FLAG_NONE, 1}
	var bytes []byte
	bytes, err = queryPacket.Send(switchIp)
	if err != nil {
		beego.Error(err)
		return err
	}
	timeAfterQuery := time.Now()
	beego.Info("Hexabus query took", timeAfterQuery.Sub(timeBeforeQuery).Seconds(), "seconds")
	beego.Trace("Query bytes", bytes)
	// Query bytes [72 88 48 67 1 0 0 0 0 1 1 1 192 215] // switching on
	// Query bytes [72 88 48 67 1 0 0 0 0 1 1 0 209 94] // switching off

	// Check packet type
	ptype, err := hexabus.PacketType(bytes)
	if err != nil {
		return err
	}
	beego.Trace("Received packet type", ptype)

	if ptype != hexabus.PTYPE_INFO {
		beego.Error("Wrong packet type - expecting hexabus InfoPacket")
		return errors.New("Failed to set switch state")
	}

	// Info packet is the only acceptable type here
	beego.Info("Received hexabus info packet")
	infoPacket := hexabus.InfoPacket{}

	// Decode info packet bytes to get the switch value as data
	err = infoPacket.Decode(bytes)
	if err != nil {
		beego.Error("Failed to decode hexabus info packet", err)
		return errors.New("Failed to set switch state")
	}

	// Expecting boolean value as data
	infoSwitchState := infoPacket.Data
	beego.Trace("Info pack switch state:", infoSwitchState)

	// The received state has to match the state written
	if switchState != infoSwitchState {
		beego.Error("Switch states do not match")
		return errors.New("Failed to set switch state")
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
