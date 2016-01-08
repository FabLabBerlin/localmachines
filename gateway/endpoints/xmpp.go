package endpoints

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"github.com/FabLabBerlin/localmachines/gateway/xmpp"
	"log"
)

type Xmpp struct {
	ns *netswitches.NetSwitches
	x  *xmpp.Xmpp
}

func NewXmpp(ns *netswitches.NetSwitches, server, user, pw string) (*Xmpp, error) {
	var err error
	x := &Xmpp{
		ns: ns,
	}
	x.x, err = xmpp.NewXmpp(server, user, pw)
	return x, err
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
	if cmd == "on" || cmd == "off" {
		return x.ns.SetOn(msg.Data.MachineId, cmd == "on")
	}
	return fmt.Errorf("invalid cmd: %v", cmd)
}
