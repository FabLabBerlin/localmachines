package endpoints

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"github.com/FabLabBerlin/localmachines/gateway/xmpp"
	"log"
	"strconv"
	"strings"
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
	if err == nil {
		x.x.Send(user, "Hello there!")
		x.x.Send(user, "Send me messages like: 2 on")
	}
	return x, err
}

func (x *Xmpp) Run() {
	log.Printf("endpoints: xmpp: Run()")
	go func() {
		for {
			select {
			case msg := <-x.x.Recv():
				if err := x.dispatch(msg); err != nil {
					log.Printf("xmpp dispatch: %v", err)
				}
			}
		}
	}()
	x.x.Run()
}

func (x *Xmpp) dispatch(msg string) (err error) {
	log.Printf("dispatch(%v)", msg)
	tmp := strings.Split(msg, " ")
	machineId, err := strconv.ParseInt(tmp[0], 10, 64)
	if err != nil {
		return fmt.Errorf("parse id: %v", err)
	}
	cmd := tmp[1]
	if cmd == "on" || cmd == "off" {
		return x.ns.SetOn(machineId, cmd == "on")
	}
	return fmt.Errorf("invalid cmd: %v", cmd)
}
