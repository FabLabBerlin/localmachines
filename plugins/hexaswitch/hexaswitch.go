package hexaswitch

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/morriswinkler/hexabus"
)

var switchAddress string = "[fafa::50:c4ff:fe04:8390]"

func On(machineId int) {
	fmt.Printf("Turning on machine with ID %d\n", machineId)
	go setSwitchState(true)
}

func Off(machineId int) {
	fmt.Printf("Turning off machine with ID %d\n", machineId)
	go setSwitchState(false)
}

func setSwitchState(switchState bool) {

	// Create write packet to switch on and off
	var wPack hexabus.WritePacket = hexabus.WritePacket{hexabus.FLAG_NONE,
		1, hexabus.DTYPE_BOOL, switchState}

	switchStateStr := "Off"
	if !switchState {
		switchStateStr = "On"
	}
	beego.Trace("Hexaswitch", switchStateStr)

	// Send packet to switch
	err := wPack.Send(switchAddress)
	if err != nil {
		fmt.Println(err)
	}
}
