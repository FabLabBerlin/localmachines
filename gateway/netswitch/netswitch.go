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
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/lib/mfi"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/FabLabBerlin/localmachines/lib/xmpp/commands"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const NETSWITCH_TYPE_MFI = "mfi"

type NetSwitch struct {
	muChInit sync.Mutex
	chSingle chan int
	On       bool `json:"-"`
	// We're using this without Beego ORM attached
	Id                  int64
	NetswitchUrlOn      string
	NetswitchUrlOff     string
	NetswitchHost       string
	NetswitchSensorPort int
	NetswitchType       string
}

func (ns *NetSwitch) SetOn(on bool) (err error) {
	if on {
		return ns.turnOn()
	} else {
		return ns.turnOff()
	}
}

func (ns *NetSwitch) turnOn() (err error) {
	log.Printf("turn on %v", ns.UrlOn())
	var resp *http.Response
	if ns.NetswitchType == NETSWITCH_TYPE_MFI {
		resp, err = http.PostForm(ns.UrlOn(), url.Values{"output": {"1"}})
	} else {
		resp, err = http.Get(ns.UrlOn())
	}
	if ns.isIgnorableAhmaError(err) {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("http: %v", err)
	}
	if resp == nil {
		log.Printf("turnOn: resp is nil!")
	} else {
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
		}
	}
	ns.On = true
	return
}

func (ns *NetSwitch) turnOff() (err error) {
	var resp *http.Response
	if ns.NetswitchType == NETSWITCH_TYPE_MFI {
		log.Printf("turn off %v", ns.UrlOn())
		resp, err = http.PostForm(ns.UrlOn(), url.Values{"output": {"0"}})
	} else {
		log.Printf("turn off %v", ns.UrlOff())
		resp, err = http.Get(ns.UrlOff())
	}
	if ns.isIgnorableAhmaError(err) {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("http: %v", err)
	}
	if resp == nil {
		log.Printf("turnOff: resp is nil!")
	} else {
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
		}
	}
	ns.On = false
	return
}

func (ns *NetSwitch) isIgnorableAhmaError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "malformed HTTP status code") &&
		strings.Contains(msg, "AhmaSwitch")
}

func (ns *NetSwitch) UrlOn() string {
	if ns.NetswitchType == NETSWITCH_TYPE_MFI {
		return "http://" + ns.NetswitchHost + "/sensors/" + strconv.Itoa(ns.NetswitchSensorPort)
	} else {
		return ns.NetswitchUrlOn
	}
}

func (ns *NetSwitch) UrlOff() string {
	if ns.NetswitchType == NETSWITCH_TYPE_MFI {
		return "http://" + ns.NetswitchHost + "/sensors/" + strconv.Itoa(ns.NetswitchSensorPort)
	} else {
		return ns.NetswitchUrlOff
	}
}

func (ns *NetSwitch) String() string {
	return fmt.Sprintf("(NetSwitch MachineId=%v On=%v)",
		ns.Id, ns.On)
}

func (ns *NetSwitch) ApplyConfig(updates chan<- string, xmppClient *xmpp.Xmpp) (err error) {
	ns.muChInit.Lock()
	if ns.chSingle == nil {
		log.Printf("make(chan int, 1)")
		ns.chSingle = make(chan int, 1)
		ns.chSingle <- 1
	} else {
		log.Printf("ns.chSingle != nil")
	}
	ns.muChInit.Unlock()
	select {
	case <-ns.chSingle:
		log.Printf("not running")
		break
	default:
		log.Println("apply config already running")
		return fmt.Errorf("apply config already running")
	}
	cfg := mfi.Config{
		Host: ns.NetswitchHost,
	}
	statusMsg := xmpp.Message{
		Remote: global.ServerJabberId,
		Data: xmpp.Data{
			Command:    commands.GATEWAY_APPLIED_CONFIG_1,
			LocationId: global.Cfg.Main.LocationId,
			MachineId:  ns.Id,
		},
	}
	if err := cfg.RunStep1WifiCredentials(); err != nil {
		ns.chSingle <- 1
		statusMsg.Data.Error = true
		statusMsg.Data.ErrorMessage = fmt.Sprintf("step 1: error getting wifi: %v", err)
		if err := xmppClient.Send(statusMsg); err != nil {
			log.Printf("xmpp command send: %v", err)
		}
		return fmt.Errorf(statusMsg.Data.ErrorMessage)
	}
	if err = xmppClient.Send(statusMsg); err != nil {
		log.Printf("xmpp command send: %v", err)
	}

	go func() {
		statusMsg := xmpp.Message{
			Remote: global.ServerJabberId,
			Data: xmpp.Data{
				Command:    commands.GATEWAY_APPLIED_CONFIG_2,
				LocationId: global.Cfg.Main.LocationId,
				MachineId:  ns.Id,
			},
		}
		if err := cfg.RunStep2PushConfig(); err != nil {
			statusMsg.Data.Error = true
			statusMsg.Data.ErrorMessage = err.Error()
			updates <- err.Error()
		}
		if err = xmppClient.Send(statusMsg); err != nil {
			log.Printf("xmpp command send: %v", err)
		}
		ns.chSingle <- 1
	}()
	return
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
