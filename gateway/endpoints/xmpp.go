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
	"time"
)

type Xmpp struct {
	ns     *netswitches.NetSwitches
	client *xmpp.Xmpp
}

func NewXmpp(ns *netswitches.NetSwitches) *Xmpp {
	x := &Xmpp{
		ns: ns,
		client: xmpp.NewXmpp(
			global.Cfg.XMPP.Server,
			global.Cfg.XMPP.User,
			global.Cfg.XMPP.Pass,
		),
	}
	x.client.Run()
	if err := x.initMachinesList(-1); err != nil {
		log.Printf("init machines list: %v", err)
	}
	for {
		select {
		case msg := <-x.client.Recv():
			log.Printf("gateway: incoming msg")
			if err := x.dispatch(msg); err != nil {
				log.Printf("error dispatching %v", msg)
			}
		}
	}
	return x
}

func (x *Xmpp) initMachinesList(retries int) (err error) {
	for i := 0; retries <= 0 || i < retries; i++ {
		if err != nil {
			log.Printf("gateway: Init: %v", err)
			log.Printf("gateway: Init: retrying in 5 seconds")
			<-time.After(5 * time.Second)
			err = nil
		}
		if err = x.ns.Load(x.client); err != nil {
			err = fmt.Errorf("netswitches load: %v", err)
			continue
		}

		if err == nil {
			log.Printf("initMachinesList: trigger Load machine list successfully.")
			break
		}
	}

	return
}

func (x *Xmpp) reinitMachinesList() (err error) {
	if err = x.initMachinesList(2); err != nil {
		return fmt.Errorf("Init: %v", err)
	}
	return
}

func (x *Xmpp) dispatch(msg xmpp.Message) (err error) {
	log.Printf("dispatch(%v)", msg.Data.Command)
	cmd := msg.Data.Command
	switch cmd {
	case "on", "off":
		resp := xmpp.Message{
			Remote: global.ServerJabberId,
			Data: xmpp.Data{
				LocationId: global.Cfg.Main.LocationId,
				MachineId:  msg.Data.MachineId,
				UserId:     msg.Data.UserId,
			},
		}
		err = x.ns.SetOn(msg.Data.MachineId, cmd == "on")
		if err == nil {
			if cmd == "on" {
				resp.Data.Command = commands.GATEWAY_SUCCESS_ON
			} else {
				resp.Data.Command = commands.GATEWAY_SUCCESS_OFF
			}
		} else {
			if cmd == "on" {
				resp.Data.Command = commands.GATEWAY_FAIL_ON
			} else {
				resp.Data.Command = commands.GATEWAY_FAIL_OFF
			}
		}

		if err = x.client.Send(resp); err != nil {
			return fmt.Errorf("xmpp command GATEWAY_ALLOWS_USERS_FROM_IP: %v", err)
		}

		return
	case commands.REINIT:
		return x.reinitMachinesList()
	case commands.APPLY_CONFIG:
		log.Printf("apply_config!!!")
		updates := make(chan string, 10)
		err := x.ns.ApplyConfig(msg.Data.MachineId, updates, x.client, msg.Data.UserId)
		log.Printf("dispatch:returning err=%v", err)
		return err
	case commands.FETCH_LOCAL_IP:
		var resp *http.Response
		resp, err = http.Get(*global.ServerPrefix + "/api/locations/my_ip")
		if err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode > 299 {
			return fmt.Errorf("status code %v", resp.StatusCode)
		}
		var buf []byte
		if buf, err = ioutil.ReadAll(resp.Body); err != nil {
			return
		}
		ipAddress := string(buf)

		if err = x.client.Send(xmpp.Message{
			Remote: global.ServerJabberId,
			Data: xmpp.Data{
				Command:    commands.GATEWAY_ALLOWS_USERS_FROM_IP,
				LocationId: global.Cfg.Main.LocationId,
				IpAddress:  ipAddress,
			},
		}); err != nil {
			return fmt.Errorf("xmpp command GATEWAY_ALLOWS_USERS_FROM_IP: %v", err)
		}
		return
	case commands.SERVER_SENDS_MACHINE_LIST:
		return x.ns.UseFromJson([]byte(msg.Data.Raw))
	}
	return fmt.Errorf("invalid cmd: %v", cmd)
}
