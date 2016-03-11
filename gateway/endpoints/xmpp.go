package endpoints

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"log"
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
				err := x.dispatch(msg)
				if err != nil {
					log.Printf("xmpp dispatch: %v", err)
				}
				response := xmpp.Message{
					Remote: msg.Remote,
					Data:   msg.Data,
				}
				response.Data.Error = err != nil
				if err := x.x.Send(response); err != nil {
					log.Printf("xmpp: failed to send response")
				}
			}
		}
	}()
	x.x.Run()
}

func (x *Xmpp) dispatch(msg xmpp.Message) (err error) {
	log.Printf("dispatch(%v)", msg)
	cmd := msg.Data.Command
	switch cmd {
	case "on", "off":
		return x.ns.SetOn(msg.Data.MachineId, cmd == "on")
	case "reinit":
		return x.reinitGateway()
	case "apply_config":
		log.Printf("apply_config!!!")
		updates := make(chan string, 10)
		err := x.ns.ApplyConfig(msg.Data.MachineId, updates)
		log.Printf("dispatch:returning err=%v", err)
		return err
	}
	return fmt.Errorf("invalid cmd: %v", cmd)
}
