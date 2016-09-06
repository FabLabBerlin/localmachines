// Gateway is a Lab "IoT" gateway server
package main

import (
	"flag"
	"github.com/FabLabBerlin/localmachines/gateway/endpoints"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"gopkg.in/gcfg.v1"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const UCI_PREFIX = "localmachines.@localmachines[0]"

func getUci(key string) (value string) {
	cmd := exec.Command("/sbin/uci", "get", UCI_PREFIX+"."+key)
	buf, err := cmd.CombinedOutput()
	value = string(buf)
	value = strings.TrimSpace(value)
	if err != nil {
		panic("get key '" + key + "': " + err.Error())
	}
	return
}

func main() {
	uci := flag.Bool("uci", false, "use UCI")
	flag.Parse()

	if *uci {
		var err error
		global.Cfg.Main.LocationId, err = strconv.ParseInt(getUci("locationid"), 10, 64)
		if err != nil {
			panic("parse locationid: " + err.Error())
		}
		global.Cfg.XMPP.User = getUci("jabberid")
		global.Cfg.XMPP.Pass = getUci("jabberpw")
		tmp := strings.Split(global.Cfg.XMPP.User, "@")
		if len(tmp) != 2 {
			panic("expected jabberid to contain exactly one '@'")
		}
		global.Cfg.XMPP.Server = tmp[1]
	} else {
		err := gcfg.ReadFileInto(&global.Cfg, "conf/gateway.conf")
		if err != nil {
			log.Printf("gcfg read file into: %v", err)
		}
	}

	netSwitches := netswitches.New()

	chTerm := make(chan os.Signal, 1)
	signal.Notify(chTerm, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range chTerm {
			log.Printf("received signal %v", sig)
			// sig is a ^C, handle it
			os.Exit(1)
		}
	}()

	go endpoints.NewXmpp(netSwitches)

	// The gateway shall run forever..
	for {
		select {}
	}
}
