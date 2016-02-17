package mocks

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/netswitch"
	modelsNetswitch "github.com/FabLabBerlin/localmachines/models/netswitch"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
)

const (
	SWITCH_ON_RESPONSE  = `{"sensors":[{"output":1,"power":4.77432537,"energy":6694.0625,"enabled":0,"current":0.04212296,"voltage":230.758893013,"powerfactor":0.491173223,"relay":1,"lock":0}],"status":"success"}`
	SWITCH_OFF_RESPONSE = `{"sensors":[{"output":0,"power":0.0,"energy":10538.4375,"enabled":1,"current":0.0,"voltage":0.0,"powerfactor":0.0,"relay":0,"lock":0}],"status":"success"}`
)

type DesiredState bool

type RelayState bool

const (
	DESIRED_ON  DesiredState = true
	DESIRED_OFF DesiredState = false
	RELAY_ON    RelayState   = true
	RELAY_OFF   RelayState   = false
)

type NetSwitch struct {
	PollRequests   int
	SwitchRequests int
	NetSwitch      *netswitch.NetSwitch
	TestServer     *httptest.Server
}

func NewNetSwitch(desired DesiredState, relay RelayState) *NetSwitch {
	mock := &NetSwitch{}
	mock.TestServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			mock.PollRequests++
		} else {
			mock.SwitchRequests++
		}
		w.WriteHeader(http.StatusOK)
		if relay == RELAY_ON {
			fmt.Fprintf(w, SWITCH_ON_RESPONSE)
		} else {
			fmt.Fprintf(w, SWITCH_OFF_RESPONSE)
		}
	}))
	url, err := url.Parse(mock.TestServer.URL)
	if err != nil {
		panic(err.Error())
	}
	mock.NetSwitch = &netswitch.NetSwitch{
		Mapping: modelsNetswitch.Mapping{
			Host:       url.Host,
			SensorPort: 1,
		},
		On: bool(desired),
	}
	return mock
}

func (ns *NetSwitch) Close() {
	ns.TestServer.Close()
}

func (ns *NetSwitch) Host() string {
	url, err := url.Parse(ns.TestServer.URL)
	if err != nil {
		panic(err.Error())
	}
	return url.Host
}

type LmApi struct {
	mu       sync.RWMutex
	mappings map[int64]modelsNetswitch.Mapping
	server   *httptest.Server
}

func NewLmApi() (api *LmApi) {
	api = &LmApi{
		mappings: make(map[int64]modelsNetswitch.Mapping),
	}
	api.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := make([]modelsNetswitch.Mapping, 0, len(api.mappings))
		for _, mapping := range api.mappings {
			state = append(state, mapping)
		}
		enc := json.NewEncoder(w)
		w.WriteHeader(http.StatusOK)
		enc.Encode(state)
	}))
	return
}

func (api *LmApi) AddMapping(mapping modelsNetswitch.Mapping) {
	api.mu.Lock()
	defer api.mu.Unlock()
	if _, ok := api.mappings[mapping.MachineId]; ok {
		panic("mapping already existing")
	}
	api.mappings[mapping.MachineId] = mapping
}

func (api *LmApi) DeleteMapping(machineId int64) {
	api.mu.Lock()
	defer api.mu.Unlock()
	delete(api.mappings, machineId)
}

func (api *LmApi) UpdateMapping(machineId int64, mapping modelsNetswitch.Mapping) {
	api.mu.Lock()
	defer api.mu.Unlock()
	if _, ok := api.mappings[machineId]; !ok {
		panic("mapping not existing yet")
	}
	api.mappings[machineId] = mapping
}

func (api *LmApi) Close() {
	api.server.Close()
}

func (api *LmApi) URL() string {
	return api.server.URL
}
