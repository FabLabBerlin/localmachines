/*
netswitches helper functions.
*/
package netswitches

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitch"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/FabLabBerlin/localmachines/lib/xmpp/commands"
	"log"
)

var (
	ErrDuplicateCombinationHostSensorPort = errors.New("Duplicate combination host + sensor port")
)

type NetSwitches struct {
	nss map[int64]*netswitch.NetSwitch
}

func New() (nss *NetSwitches) {
	nss = &NetSwitches{}
	return
}

// Load netswitches from EASY LAB API.  client should be logged in.
func (nss *NetSwitches) Load(xmppClient *xmpp.Xmpp) (err error) {
	if err = xmppClient.Send(xmpp.Message{
		Remote: global.ServerJabberId,
		Data: xmpp.Data{
			Command:    commands.GATEWAY_REQUESTS_MACHINE_LIST,
			LocationId: global.Cfg.Main.LocationId,
		},
	}); err != nil {
		return fmt.Errorf("xmpp command GATEWAY_REQUESTS_MACHINE_LIST: %v", err)
	}
	log.Printf("netswitches: %v", nss.nss)
	return
}

// fetch netswitch data from EASY LAB API.
//
// Each type NetSwitch runs its own dispatch loop.  Make sure no additional
// loop is started.
func (nss *NetSwitches) UseFromJson(raw []byte) (err error) {
	all := []netswitch.NetSwitch{}
	if err := json.Unmarshal(raw, &all); err != nil {
		return fmt.Errorf("json unmarshal: %v", err)
	}
	mappings := make([]netswitch.NetSwitch, 0, len(all))
	for _, mapping := range all {
		mappings = append(mappings, mapping)
	}
	if nss.nss == nil {
		nss.nss = make(map[int64]*netswitch.NetSwitch)
	}
	for _, mapping := range mappings {
		log.Printf("adding mapping with id %v", mapping.Id)
		log.Printf("host: %v", mapping.NetswitchHost)
		if _, exists := nss.nss[mapping.Id]; !exists {
			nss.nss[mapping.Id] = &netswitch.NetSwitch{
				Id: mapping.Id,
			}
		}
		existing := nss.nss[mapping.Id]
		log.Printf("exists")
		existing.NetswitchUrlOn = mapping.NetswitchUrlOn
		existing.NetswitchUrlOff = mapping.NetswitchUrlOff
		existing.NetswitchHost = mapping.NetswitchHost
		existing.NetswitchSensorPort = mapping.NetswitchSensorPort
		existing.NetswitchType = mapping.NetswitchType
	}
	log.Printf("now mappings nss[2]=%v", nss.nss[2])
	// Removed unused IDs
	unusedIDs := make([]int64, 0, 3)
	for _, ns := range nss.nss {
		foundInMappings := false
		for _, mapping := range mappings {
			if mapping.Id == ns.Id {
				foundInMappings = true
			}
		}
		if !foundInMappings {
			unusedIDs = append(unusedIDs, ns.Id)
		}
	}
	log.Printf("unusedIds=%v", unusedIDs)
	for _, unusedID := range unusedIDs {
		delete(nss.nss, unusedID)
	}
	// Make sure there are no duplicate combinations Netswitch Host + SensorPort
	hostsSensorPorts := make(map[string]bool)
	for _, ns := range nss.nss {
		if ns.NetswitchType != "" && len(ns.NetswitchHost) > 0 {
			key := fmt.Sprintf("%v -> %v", ns.NetswitchHost, ns.NetswitchSensorPort)
			if _, ok := hostsSensorPorts[key]; ok {
				log.Printf("duplicate combination host sensor port: %v",
					ns.NetswitchHost, ns.NetswitchSensorPort)
				return ErrDuplicateCombinationHostSensorPort
			}
			hostsSensorPorts[key] = true
		}
	}
	return
}

func (nss *NetSwitches) SetOn(machineId int64, on bool) (err error) {
	for retries := 0; retries < global.MAX_SYNC_RETRIES; retries++ {
		if err = nss.setOn(machineId, on); err == nil {
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

func (nss *NetSwitches) ApplyConfig(machineId int64, updates chan<- string, xmppClient *xmpp.Xmpp) (err error) {
	ns, ok := nss.nss[machineId]
	if !ok {
		return fmt.Errorf("no netswitch for machine id %v present",
			machineId)
	}
	return ns.ApplyConfig(updates, xmppClient)
}
