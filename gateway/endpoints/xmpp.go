package endpoints

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/FabLabBerlin/localmachines/lib/xmpp/commands"
	"io/ioutil"
	"log"
	"net/http"
)

type Xmpp struct {
	ns            *netswitches.NetSwitches
	x             *xmpp.Xmpp
	reinitGateway func() error
}

func NewXmpp(ns *netswitches.NetSwitches, reinitGateway func() error) *Xmpp {
	x := &Xmpp{
		ns:            ns,
		reinitGateway: reinitGateway,
	}
	x.x = xmpp.NewXmpp(global.Cfg.XMPP.Server, global.Cfg.XMPP.User, global.Cfg.XMPP.Pass)
	return x
}

func (x *Xmpp) Run() {
	log.Printf("endpoints: xmpp: Run()")
	go func() {
		for {
			select {
			case msg := <-x.x.Recv():
				ipAddress, err := x.dispatch(msg)
				if err != nil {
					log.Printf("xmpp dispatch: %v", err)
				}
				response := xmpp.Message{
					Remote: msg.Remote,
					Data:   msg.Data,
				}
				response.Data.IpAddress = ipAddress
				response.Data.Error = err != nil
				if err := x.x.Send(response); err != nil {
					log.Printf("xmpp: failed to send response")
				}
			}
		}
	}()
	x.x.Run()
}

func (x *Xmpp) dispatch(msg xmpp.Message) (ipAddress string, err error) {
	log.Printf("dispatch(%v)", msg)
	cmd := msg.Data.Command
	switch cmd {
	case "on", "off":
		return "", x.ns.SetOn(msg.Data.MachineId, cmd == "on")
	case commands.REINIT:
		return "", x.reinitGateway()
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
