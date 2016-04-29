// Gateway is a Lab "IoT" gateway server
package main

import (
	"encoding/json"
	"flag"
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
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const UCI_PREFIX = "localmachines.@localmachines[0]"

var netSwitches *netswitches.NetSwitches

type LoginResp struct {
	Status string
	UserId int64
}

func (resp *LoginResp) ok() bool {
	return resp.Status == "ok"
}

func Login(client *http.Client, user, key string) (err error) {
	locationId := strconv.FormatInt(global.Cfg.Main.LocationId, 10)
	uri := global.Cfg.API.Url + "/users/login?location=" + locationId
	resp, err := client.PostForm(uri,
		url.Values{"username": {user}, "password": {key}})
	if err != nil {
		return fmt.Errorf("POST login: %v", err)
	}
	defer resp.Body.Close()
	if code := resp.StatusCode; code > 299 {
		return fmt.Errorf("unexpected status code: %v", code)
	}
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
	if err = Init(2); err != nil {
		return fmt.Errorf("Init: %v", err)
	}
	return
}

func getUci(key string) (value string) {
	cmd := exec.Command("/sbin/uci", "get", UCI_PREFIX + "." + key)
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
		global.Cfg.API.Id = getUci("apiid")
		global.Cfg.API.Key = getUci("apikey")
		global.Cfg.API.Url = "https://easylab.io/api"
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

	go endpoints.NewHttp()

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

	endpoints.NewXmpp(netSwitches, Reinit)

	// The gateway shall run forever..
	for {
		select{}
	}
}
