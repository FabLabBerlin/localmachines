/*
netswitch control.

The NetSwitches have Relays that must be in sync with our system.  For example
when a NetSwitch is plugged in, the relay must go into the correct position.
The mfi Switches are by default switched on and need to turn off after being
plugged in.
*/
package netswitch

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type syncCommand struct {
	SetOn *bool
	Error chan error
}

type NetSwitch struct {
	machine.Machine
	On     bool
	syncCh chan syncCommand
}

func (ns *NetSwitch) assertLoopRunning() {
	if ns.syncCh == nil {
		ns.syncCh = make(chan syncCommand, 1)
		go ns.loop()
	}
}

func (ns *NetSwitch) loop() {
	for cmd := range ns.syncCh {
		cmd.Error <- ns.sync(cmd)
	}
}

func (ns *NetSwitch) Close() {
	log.Printf("NetSwitch#Close")
	if ns.syncCh != nil {
		close(ns.syncCh)
	}
}

func (ns *NetSwitch) sync(cmd syncCommand) (err error) {
	log.Printf("urlStatus = %v", ns.Url())
	if ns.NetswitchSensorPort < 1 {
		return fmt.Errorf("NetSwitch switch port for Machine %v invalid", ns.Id)
	}
	client := http.Client{
		Timeout: global.STATE_SYNC_TIMEOUT,
	}
	resp, err := client.Get(ns.Url())
	if err != nil {
		return fmt.Errorf("http get url status: %v", err)
	}
	defer resp.Body.Close()
	if code := resp.StatusCode; code < 200 || code > 299 {
		return fmt.Errorf("unexpected status code: %v", code)
	}
	mfi := MfiSwitch{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&mfi); err != nil {
		return fmt.Errorf("json decode:", err)
	}
	onDesired := ns.On
	if cmd.SetOn != nil {
		onDesired = *cmd.SetOn
	}
	if mfi.On() != onDesired {
		log.Printf("State for Machine %v must get synchronized", ns.Id)
		if onDesired {
			if err = ns.turnOn(); err != nil {
				return fmt.Errorf("turn on: %v", err)
			}
		} else {
			if err = ns.turnOff(); err != nil {
				return fmt.Errorf("turn off: %v", err)
			}
		}
	}
	return
}

func (ns *NetSwitch) SetOn(on bool) (err error) {
	ns.assertLoopRunning()
	cmd := syncCommand{
		SetOn: &on,
		Error: make(chan error),
	}
	ns.syncCh <- cmd
	return <-cmd.Error
}

func (ns *NetSwitch) Sync() (err error) {
	ns.assertLoopRunning()
	cmd := syncCommand{
		Error: make(chan error),
	}
	ns.syncCh <- cmd
	return <-cmd.Error
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
