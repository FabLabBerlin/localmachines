package gatewayNetswitchTest

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/netswitch"
	modelsNetswitch "github.com/FabLabBerlin/localmachines/models/netswitch"
	"net/http"
	"net/http/httptest"
	"net/url"
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

type MockNetSwitch struct {
	UrlCalled  int
	NetSwitch  *netswitch.NetSwitch
	TestServer *httptest.Server
}

func NewMockNetSwitch(desired DesiredState, relay RelayState) *MockNetSwitch {
	mock := &MockNetSwitch{}
	mock.TestServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.UrlCalled++
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
