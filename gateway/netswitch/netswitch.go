/*
netswitch control.

The NetSwitches have Relays that must be in sync with our system.  For example
when a NetSwitch is plugged in, the relay must go into the correct position.
The mfi Switches are by default switched on and need to turn off after being
plugged in.
*/
package netswitch

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type NetSwitch struct {
	machine.Machine
	On bool
}

func (ns *NetSwitch) SetOn(on bool) (err error) {
	if on {
		return ns.turnOn()
	} else {
		return ns.turnOff()
	}
}

func (ns *NetSwitch) turnOn() (err error) {
	log.Printf("turn on %v", ns.Url())
	resp, err := http.PostForm(ns.Url(), url.Values{"output": {"1"}})
	if err != nil {
		return fmt.Errorf("client get: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	ns.On = true
	return
}

func (ns *NetSwitch) turnOff() (err error) {
	log.Printf("turn off %v", ns.Url())
	resp, err := http.PostForm(ns.Url(), url.Values{"output": {"0"}})
	if err != nil {
		return fmt.Errorf("client get: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	ns.On = false
	return
}

func (ns *NetSwitch) Url() string {
	return "http://" + ns.NetswitchHost + "/sensors/" + strconv.Itoa(ns.NetswitchSensorPort)
}

func (ns *NetSwitch) String() string {
	return fmt.Sprintf("(NetSwitch MachineId=%v On=%v)",
		ns.Id, ns.On)
}

//{"sensors":[{"output":1,"power":0.0,"energy":0.0,"enabled":0,"current":0.0,"voltage":233.546874046,"powerfactor":0.0,"relay":1,"lock":0}],"status":"success"}

type MfiSwitch struct {
	Sensors []MfiSensor `json:"sensors"`
	Status  string      `json:"status"`
}

func (swi *MfiSwitch) On() bool {
	relay := swi.Sensors[0].Relay
	switch relay {
	case 0:
		return false
		break
	case 1:
		return true
		break
	}
	log.Fatalf("unknown relay status %v, terminating", relay)
	return false
}

type MfiSensor struct {
	Output      int     `json:"output"`
	Power       float64 `json:"power"`
	Energy      float64 `json:"energy"`
	Enabled     float64 `json:"enabled"`
	Current     float64 `json:"current"`
	Voltage     float64 `json:"voltage"`
	PowerFactor float64 `json:"powerfactor"`
	Relay       int     `json:"relay"`
	Lock        int     `json:"lock"`
}
