package main

import (
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
)

type LoginResp struct {
	Status string
	UserId int64
}

func (resp *LoginResp) ok() bool {
	return resp.Status == "ok"
}

var client *http.Client
var netSwitches map[int64]*NetSwitch

type NetSwitch struct {
	Id        int64
	MachineId int64
	UrlOn     string
	UrlOff    string
	On        bool
}

func Login(apiUrl, user, key string) (err error) {
	resp, err := client.PostForm(apiUrl+"/users/login",
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

func Fetch(apiUrl string) (err error) {
	resp, err := client.Get(apiUrl + "/netswitch")
	if err != nil {
		return fmt.Errorf("GET: %v", err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	nss := []NetSwitch{}
	if err := dec.Decode(&nss); err != nil {
		return fmt.Errorf("json decode: %v", err)
	}
	log.Printf("netswitches: %v", nss)
	netSwitches = make(map[int64]*NetSwitch)
	for _, ns := range nss {
		netSwitches[ns.MachineId] = &ns
	}
	return
}

func loadState(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var switchStates []NetSwitch
	if err := dec.Decode(&switchStates); err != nil {
		return fmt.Errorf("json decode: %v", err)
	}
	for _, switchState := range switchStates {
		mid := switchState.MachineId
		if netswitch, ok := netSwitches[mid]; ok {
			netswitch.On = switchState.On
		} else {
			log.Printf("netswitch for machine id %v doesn't exist anymore", mid)
		}
	}
	return
}

func saveState(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	switchStates := make([]NetSwitch, 0, len(netSwitches))
	for _, ns := range netSwitches {
		switchStates = append(switchStates, *ns)
	}
	return enc.Encode(switchStates)
}

func SaveState(filename string) {
	log.Printf("Saving state...")
	if err := saveState(filename); err != nil {
		log.Printf("failed saving state: %v", err)
	}
}

func Init(apiUrl, user, key, stateFile string) (err error) {

	client = &http.Client{}
	if client.Jar, err = cookiejar.New(nil); err != nil {
		return
	}
	if err := Login(apiUrl, user, key); err != nil {
		return fmt.Errorf("login: %v", err)
	}
	if err := Fetch(apiUrl); err != nil {
		return fmt.Errorf("fetch: %v", err)
	}
	if err := loadState(stateFile); err != nil {
		return fmt.Errorf("load state: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("received signal %v", sig)
			// sig is a ^C, handle it
			SaveState(stateFile)
			os.Exit(1)
		}
	}()

	return
}

func runCommand(w http.ResponseWriter, r *http.Request) (err error) {
	tmp := strings.Split(r.URL.Path, "/")
	idStr := tmp[len(tmp)-2]
	cmd := tmp[len(tmp)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parse id: %v", err)
	}
	ns, ok := netSwitches[id]
	if !ok {
		return fmt.Errorf("there's no netswitch for machine id %v", id)
	}
	switch cmd {
	case "on":
		resp, err := client.Get(ns.UrlOn)
		if err != nil {
			return fmt.Errorf("client get: %v", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
		}
		ns.On = true
		break
	case "off":
		resp, err := client.Get(ns.UrlOff)
		if err != nil {
			return fmt.Errorf("client get: %v", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
		}
		ns.On = false
		break
	default:
		return fmt.Errorf("unknown command '%v'", cmd)
	}
	log.Printf("id: %v", id)
	log.Printf("cmd: %v", cmd)
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
	apiUrl := flag.String("apiUrl", "http://localhost:8080/api", "Url of the fabsmith api (http or https)")
	user := flag.String("id", "user", "id")
	key := flag.String("key", "user", "key")
	stateFile := flag.String("stateFile", "state.json", "switches are stateful but they loose state on reset")
	flag.Parse()
	if err := Init(*apiUrl, *user, *key, *stateFile); err != nil {
		log.Fatalf("Init: %v", err)
	}

	http.HandleFunc("/machines/", RunCommand)

	if err := http.ListenAndServe(":7070", nil); err != nil {
		SaveState(*stateFile)
	}
}
