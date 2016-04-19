package endpoints

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/FabLabBerlin/localmachines/lib/xmpp/commands"
	"github.com/FabLabBerlin/localmachines/lib/xmpp/request_response"
	"io/ioutil"
	"log"
	"net/http"
)

type Xmpp struct {
	ns            *netswitches.NetSwitches
	dispatcher    *request_response.Dispatcher
	reinitGateway func(payloadJson string) error
}

func NewXmpp(ns *netswitches.NetSwitches, reinitGateway func(payloadJson string) error) *Xmpp {
	x := &Xmpp{
		ns:            ns,
		reinitGateway: reinitGateway,
	}
	x.dispatcher = request_response.NewDispatcher(global.Cfg.XMPP.Server, global.Cfg.XMPP.User, global.Cfg.XMPP.Pass, x.dispatch)
	return x
}

func (x *Xmpp) dispatch(msg xmpp.Message) (ipAddress string, err error) {
	log.Printf("dispatch(%v)", msg)
	cmd := msg.Data.Command
	switch cmd {
	case "on", "off":
		return "", x.ns.SetOn(msg.Data.MachineId, cmd == "on")
	case commands.REINIT:
		return "", x.reinitGateway(msg.Data.Payload)
	case commands.APPLY_CONFIG:
		log.Printf("apply_config!!!")
		updates := make(chan string, 10)
		err := x.ns.ApplyConfig(msg.Data.MachineId, updates)
		log.Printf("dispatch:returning err=%v", err)
		return "", err
	case commands.FETCH_LOCAL_IP:
		var resp *http.Response
		resp, err = http.Get(global.Cfg.API.Url + "/locations/my_ip")
		if err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode > 299 {
			return "", fmt.Errorf("status code %v", resp.StatusCode)
		}
		var buf []byte
		if buf, err = ioutil.ReadAll(resp.Body); err != nil {
			return
		}
		ipAddress = string(buf)
		return
	}
	return "", fmt.Errorf("invalid cmd: %v", cmd)
}

func (x *Xmpp) RequestReinit() (err error) {
	log.Printf("xmpp endpoint: RequestReinit()")
	locId := global.Cfg.Main.LocationId
	_, err = x.dispatcher.SendXmppCommand(locId, global.Cfg.XMPP.MainId, commands.REQUEST_REINIT, 0, "")
	return
}
