package netswitches

import (
	"../global"
	"../netswitch"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type NetSwitches struct {
	syncCh chan syncCommand
	nss    map[int64]*netswitch.NetSwitch
}

type syncCommand struct {
	MachineId *int64
	Error     chan error
}

func New() (nss *NetSwitches) {
	nss = &NetSwitches{
		syncCh: make(chan syncCommand, 10),
	}
	go nss.dispatchLoop()
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
	nss.nss = make(map[int64]*netswitch.NetSwitch)
	for _, ns := range list {
		nss.nss[ns.MachineId] = ns
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
			netswitch.SetOn(switchState.On())
			nss.nss[mid] = netswitch
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

func (nss *NetSwitches) SetOn(machineId int64, on bool) {
	ns := nss.nss[machineId]
	ns.SetOn(on)
	nss.nss[machineId] = ns
}

func (nss *NetSwitches) Sync(machineId *int64) (err error) {
	cmd := syncCommand{
		MachineId: machineId,
		Error:     make(chan error),
	}
	nss.syncCh <- cmd
	if err := <-cmd.Error; err != nil {
		return fmt.Errorf("cmd dispatch: %v", err)
	}
	return
}

func (nss *NetSwitches) dispatch(cmd syncCommand) (err error) {
	wg := sync.WaitGroup{}
	var errs error
	for _, ns := range nss.nss {
		if mid := cmd.MachineId; mid != nil && *mid != ns.MachineId {
			continue
		}
		wg.Add(1)
		go func(ns *netswitch.NetSwitch) {
			if err := ns.Sync(); err != nil {
				if errs == nil {
					errs = err
				} else {
					errs = fmt.Errorf("%v; %v", errs, err)
				}
			}
			wg.Done()
		}(ns)
	}
	wg.Wait()
	if errs != nil {
		return fmt.Errorf("state watch: %v", errs)
	}
	return
}

func (nss *NetSwitches) dispatchLoop() {
	for {
		select {
		case cmd := <-nss.syncCh:
			cmd.Error <- nss.dispatch(cmd)
		}
	}
}
