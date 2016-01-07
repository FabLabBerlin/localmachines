package netswitches

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitch"
	"log"
	"net/http"
	"os"
)

type NetSwitches struct {
	nss map[int64]*netswitch.NetSwitch
}

func New() (nss *NetSwitches) {
	nss = &NetSwitches{}
	return
}

// Load netswitches from EASY LAB API and local state json.  client should
// be logged in.
func (nss *NetSwitches) Load(client *http.Client) (err error) {
	if err = nss.fetch(client); err != nil {
		return fmt.Errorf("fetch: %v", err)
	}
	if err = nss.loadOnOff(); err != nil {
		return fmt.Errorf("load on off: %v", err)
	}
	log.Printf("netswitches: %v", nss.nss)
	return
}

// fetch netswitch data from EASY LAB API.
//
// Each type NetSwitch runs its own dispatch loop.  Make sure no additional
// loop is started.  TODO: get rid of removed NetSwitches
func (nss *NetSwitches) fetch(client *http.Client) (err error) {
	resp, err := client.Get(global.ApiUrl + "/netswitch")
	if err != nil {
		return fmt.Errorf("GET: %v", err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	list := []*netswitch.NetSwitch{}
	if err := dec.Decode(&list); err != nil {
		return fmt.Errorf("json decode: %v", err)
	}
	if nss.nss == nil {
		nss.nss = make(map[int64]*netswitch.NetSwitch)
	}
	for _, ns := range list {
		if existing, exists := nss.nss[ns.MachineId]; exists {
			existing.Id = ns.Id
			existing.UrlOn = ns.UrlOn
			existing.UrlOff = ns.UrlOff
		} else {
			nss.nss[ns.MachineId] = ns
		}
	}
	return
}

func (nss *NetSwitches) loadOnOff() (err error) {
	f, err := os.Open(global.StateFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var switchStates []*netswitch.NetSwitch
	if err := dec.Decode(&switchStates); err != nil {
		return fmt.Errorf("json decode: %v", err)
	}
	for _, switchState := range switchStates {
		mid := switchState.MachineId
		if netswitch, ok := nss.nss[mid]; ok {
			netswitch.On = switchState.On
		} else {
			log.Printf("netswitch for machine id %v doesn't exist anymore", mid)
		}
	}
	return
}

func (nss *NetSwitches) save(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	switchStates := make([]*netswitch.NetSwitch, 0, len(nss.nss))
	for _, ns := range nss.nss {
		switchStates = append(switchStates, ns)
	}
	return enc.Encode(switchStates)
}

func (nss *NetSwitches) Save() {
	log.Printf("Saving state...")
	if err := nss.save(global.StateFilename); err != nil {
		log.Printf("failed saving state: %v", err)
	}
}

func (nss *NetSwitches) SetOn(machineId int64, on bool) (err error) {
	for retries := 0; retries < global.MAX_SYNC_RETRIES; retries++ {
		if err = nss.SetOn(machineId, on); err == nil {
			if retries > 0 {
				log.Printf("Synchronized netswitch after %v tries", retries+1)
			}
			return
		}
	}
	return
}

func (nss *NetSwitches) setOn(machineId int64, on bool) (err error) {
	ns, ok := nss.nss[machineId]
	if !ok {
		return fmt.Errorf("no netswitch for machine id %v present",
			machineId)
	}
	return ns.SetOn(on)
}

func (nss *NetSwitches) Sync(machineId int64) (err error) {
	ns, ok := nss.nss[machineId]
	if !ok {
		return fmt.Errorf("no netswitch for machine id %v present",
			machineId)
	}
	return ns.Sync()
}

func (nss *NetSwitches) SyncAll() error {
	var errs error
	for _, ns := range nss.nss {
		if err := ns.Sync(); err != nil {
			if errs == nil {
				errs = err
			} else {
				errs = fmt.Errorf("%v; %v", errs, err)
			}
		}
	}
	return errs
}
