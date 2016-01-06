package netswitch

import (
	"../global"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type NetSwitch struct {
	Id        int64
	MachineId int64
	UrlOn     string
	UrlOff    string
	On        bool
}

func (ns NetSwitch) Sync() (err error) {
	urlStatus := ns.UrlStatus()
	log.Printf("urlStatus = %v", urlStatus)
	if urlStatus == "" {
		return fmt.Errorf("NetSwitch status url for Machine %v empty", ns.MachineId)
	}
	client := http.Client{
		Timeout: global.STATE_SYNC_TIMEOUT,
	}
	resp, err := client.Get(urlStatus)
	if err != nil {
		return fmt.Errorf("http get url status: %v", err)
	}
	defer resp.Body.Close()
	mfi := MfiSwitch{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&mfi); err != nil {
		return fmt.Errorf("json decode:", err)
	}
	if mfi.On() != ns.On {
		log.Printf("State for Machine %v must get synchronized", ns.MachineId)
		if ns.On {
			if err = ns.TurnOn(); err != nil {
				return fmt.Errorf("turn on: %v", err)
			}
		} else {
			if err = ns.TurnOff(); err != nil {
				return fmt.Errorf("turn off: %v", err)
			}
		}
	}
	return
}

func (ns NetSwitch) TurnOn() (err error) {
	log.Printf("turn on %v", ns.UrlOn)
	resp, err := http.Get(ns.UrlOn)
	if err != nil {
		return fmt.Errorf("client get: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return
}

func (ns NetSwitch) TurnOff() (err error) {
	log.Printf("turn off %v", ns.UrlOff)
	resp, err := http.Get(ns.UrlOff)
	if err != nil {
		return fmt.Errorf("client get: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return
}

// UrlStatus is a really dirty function.  It's just here for the proof of concept.
func (ns NetSwitch) UrlStatus() string {
	tmp := strings.Split(ns.UrlOn, "//")
	host := strings.Split(tmp[1], "/")[0]
	return "http://" + host + "/sensors/1"
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
