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

var (
	cfg         Config
	netSwitches *netswitches.NetSwitches
)

type Config struct {
	Main struct {
		StateFile string
	}
	API struct {
		Id  string
		Key string
		Url string
	}
	XMPP struct {
		Server string
		User   string
		Pass   string
	}
}

type LoginResp struct {
	Status string
	UserId int64
}

func (resp *LoginResp) ok() bool {
	return resp.Status == "ok"
}

func Login(client *http.Client, user, key string) (err error) {
	resp, err := client.PostForm(global.ApiUrl+"/users/login",
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

func Init(user, key string) (err error) {
	client := &http.Client{}
	if client.Jar, err = cookiejar.New(nil); err != nil {
		return
	}
	if err := Login(client, user, key); err != nil {
		return fmt.Errorf("login: %v", err)
	}
	if err = netSwitches.Load(client); err != nil {
		return fmt.Errorf("netswitches load: %v", err)
	}

	return
}

func main() {
	err := gcfg.ReadFileInto(&cfg, "conf/gateway.conf")
	if err != nil {
		log.Fatalf("gcfg read file into: %v", err)
	}
	global.ApiUrl = cfg.API.Url
	global.StateFilename = cfg.Main.StateFile

	netSwitches = netswitches.New()

	if err := Init(cfg.API.Id, cfg.API.Key); err != nil {
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
			netSwitches.Save()
			if err := Init(cfg.API.Id, cfg.API.Key); err != nil {
				log.Fatalf("Init: %v", err)
			}
		}
	}()

	go PingLoop()

	xmpp, err := endpoints.NewXmpp(netSwitches, cfg.XMPP.Server, cfg.XMPP.User, cfg.XMPP.Pass)
	if err != nil {
		log.Fatalf("xmpp: %v", err)
	}
	xmpp.Run()

	httpServer := endpoints.NewHttpServer(netSwitches)
	httpServer.Run()

}
