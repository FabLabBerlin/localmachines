// Gateway is a Lab "IoT" gateway server
package main

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/endpoints"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/FabLabBerlin/localmachines/gateway/netswitches"
	"gopkg.in/gcfg.v1"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var netSwitches *netswitches.NetSwitches

type LoginResp struct {
	Status string
	UserId int64
}

func (resp *LoginResp) ok() bool {
	return resp.Status == "ok"
}

func Login(client *http.Client, user, key string) (err error) {
	resp, err := client.PostForm(global.Cfg.API.Url+"/users/login",
		url.Values{"username": {user}, "password": {key}})
	if err != nil {
		return fmt.Errorf("POST login: %v", err)
	}
	defer resp.Body.Close()
	loginResp := LoginResp{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&loginResp); err != nil {
		return fmt.Errorf("json decode: %v", err)
	}
	if !loginResp.ok() {
		return fmt.Errorf("login failed: %v", loginResp.Status)
	}
	log.Printf("Logged in with user id %v", loginResp.UserId)
	return
}

func PingLoop() {
	for {
		select {
		case <-time.After(global.STATE_SYNC_PERIOD):
			if err := netSwitches.SyncAll(); err != nil {
				log.Printf("state watch err: %v", err)
			}
		}
	}
}

func Init(retries int) (err error) {
	user := global.Cfg.API.Id
	key := global.Cfg.API.Key

	for i := 0; retries <= 0 || i < retries; i++ {
		client := &http.Client{}
		if err != nil {
			log.Printf("gateway: Init: %v", err)
			log.Printf("gateway: Init: retrying in 5 seconds")
			<-time.After(5 * time.Second)
			err = nil
		}
		if client.Jar, err = cookiejar.New(nil); err != nil {
			err = fmt.Errorf("cookie jar: %v", err)
			continue
		}
		if err = Login(client, user, key); err != nil {
			err = fmt.Errorf("login: %v", err)
			continue
		}
		if err = netSwitches.Load(client); err != nil {
			err = fmt.Errorf("netswitches load: %v", err)
			continue
		}
		if err == nil {
			break
		}
	}

	return
}

func Reinit() (err error) {
	netSwitches.Save()
	if err = Init(2); err != nil {
		return fmt.Errorf("Init: %v", err)
	}
	return
}

func main() {
	err := gcfg.ReadFileInto(&global.Cfg, "conf/gateway.conf")
	if err != nil {
		log.Fatalf("gcfg read file into: %v", err)
	}

	netSwitches = netswitches.New()

	if err := Init(-1); err != nil {
		log.Fatalf("Init: %v", err)
	}

	chTerm := make(chan os.Signal, 1)
	signal.Notify(chTerm, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range chTerm {
			log.Printf("received signal %v", sig)
			// sig is a ^C, handle it
			netSwitches.Save()
			os.Exit(1)
		}
	}()

	chHup := make(chan os.Signal, 1)
	signal.Notify(chHup, syscall.SIGHUP)
	go func() {
		for sig := range chHup {
			log.Printf("received signal %v", sig)
			Reinit()
		}
	}()

	go PingLoop()

	xmpp := endpoints.NewXmpp(netSwitches, Reinit)
	xmpp.Run()

	httpServer := endpoints.NewHttpServer(netSwitches)
	httpServer.Run()

}
