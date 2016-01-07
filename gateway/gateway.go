package main

import (
	"./global"
	"./netswitches"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type CommandType string

const (
	CMD_ON  = "on"
	CMD_OFF = "off"
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

func runCommand(w http.ResponseWriter, r *http.Request) (err error) {
	tmp := strings.Split(r.URL.Path, "/")
	idStr := tmp[len(tmp)-2]
	cmdStr := tmp[len(tmp)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parse id: %v", err)
	}
	log.Printf("id: %v", id)
	log.Printf("cmd: %v", cmdStr)

	switch cmdStr {
	case CMD_ON, CMD_OFF:
		for retries := 0; retries < global.MAX_SYNC_RETRIES; retries++ {
			if err = netSwitches.SetOn(id, cmdStr == CMD_ON); err == nil {
				if retries > 0 {
					log.Printf("Synchronized netswitch after %v tries", retries+1)
				}
				return
			}
		}
		break
	default:
		return fmt.Errorf("unknown command '%v'", cmdStr)
	}
	return
}

func RunCommand(w http.ResponseWriter, r *http.Request) {
	if err := runCommand(w, r); err == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %v", err)
		log.Printf("run command: %v", err)
	}
}

func main() {
	global.ApiUrl = *flag.String("apiUrl", "http://localhost:8080/api", "Url of the fabsmith api (http or https)")
	user := flag.String("id", "user", "id")
	key := flag.String("key", "user", "key")
	global.StateFilename = *flag.String("stateFile", "state.json", "switches are stateful but they loose state on reset")
	flag.Parse()

	netSwitches = netswitches.New()

	if err := Init(*user, *key); err != nil {
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
			if err := Init(*user, *key); err != nil {
				log.Fatalf("Init: %v", err)
			}
		}
	}()

	go PingLoop()

	http.HandleFunc("/machines/", RunCommand)

	if err := http.ListenAndServe(":7070", nil); err != nil {
		netSwitches.Save()
	}
}
