// Gateway is a Lab "IoT" gateway server
package main

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/endpoints"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"gopkg.in/gcfg.v1"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	endpoint    *endpoints.Xmpp
	netSwitches *netswitches.NetSwitches
)

func Init(retries int) (err error) {
	for i := 0; retries <= 0 || i < retries; i++ {
		if err = endpoint.RequestReinit(); err != nil {
			err = fmt.Errorf("request reinit: %v", err)
			log.Printf("%v", err)
			log.Printf("gateway: Init: retrying in 5 seconds")
			<-time.After(5 * time.Second)
			err = nil
		} else {
			break
		}
	}

	return
}

func Reinit(payloadJson string) (err error) {
	if err = Init(2); err != nil {
		return fmt.Errorf("Init: %v", err)
	}
	if err = netSwitches.Load([]byte(payloadJson)); err != nil {
		return fmt.Errorf("load netswitches: %v", err)
	}
	return
}

func main() {
	err := gcfg.ReadFileInto(&global.Cfg, "conf/gateway.conf")
	if err != nil {
		log.Fatalf("gcfg read file into: %v", err)
	}

	netSwitches = netswitches.New()

	chTerm := make(chan os.Signal, 1)
	signal.Notify(chTerm, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range chTerm {
			log.Printf("received signal %v", sig)
			// sig is a ^C, handle it
			os.Exit(1)
		}
	}()

	endpoint = endpoints.NewXmpp(netSwitches, Reinit)

	if err := Init(-1); err != nil {
		log.Fatalf("Init: %v", err)
	}

	// The gateway shall run forever..
	for {
		select{}
	}
}
