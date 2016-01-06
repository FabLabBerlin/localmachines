package main

import (
	"./global"
	"./netswitch"
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
	"sync"
	"time"
)

type CommandType string

const (
	CMD_ON         = "on"
	CMD_OFF        = "off"
	CMD_STATE_SYNC = "state_sync"
)

var apiUrl string
var state *State

// Command channel has a high water mark of 10 commands in the queue.
var ch = make(chan Command, 10)

type Command struct {
	CommandType
	MachineId *int64
	Error     chan error
}

type LoginResp struct {
	Status string
	UserId int64
}

func (resp *LoginResp) ok() bool {
	return resp.Status == "ok"
}

func Login(client *http.Client, user, key string) (err error) {
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

type State struct {
	netSwitches map[int64]netswitch.NetSwitch
	filename    string
}

func LoadState(filename string, client *http.Client) (s *State, err error) {
	s = &State{
		filename: filename,
	}
	if err = s.fetchSwitches(client); err != nil {
		return nil, fmt.Errorf("fetch switches: %v", err)
	}
	if err = s.loadOnOff(); err != nil {
		return nil, fmt.Errorf("load on off: %v", err)
	}
	return
}

func (s *State) fetchSwitches(client *http.Client) (err error) {
	resp, err := client.Get(apiUrl + "/netswitch")
	if err != nil {
		return fmt.Errorf("GET: %v", err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	nss := []netswitch.NetSwitch{}
	if err := dec.Decode(&nss); err != nil {
		return fmt.Errorf("json decode: %v", err)
	}
	s.netSwitches = make(map[int64]netswitch.NetSwitch)
	for _, ns := range nss {
		s.netSwitches[ns.MachineId] = ns
	}
	return
}

func (s *State) loadOnOff() (err error) {
	f, err := os.Open(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var switchStates []netswitch.NetSwitch
	if err := dec.Decode(&switchStates); err != nil {
		return fmt.Errorf("json decode: %v", err)
	}
	for _, switchState := range switchStates {
		mid := switchState.MachineId
		if netswitch, ok := s.netSwitches[mid]; ok {
			netswitch.On = switchState.On
			s.netSwitches[mid] = netswitch
		} else {
			log.Printf("netswitch for machine id %v doesn't exist anymore", mid)
		}
	}
	return
}

func (s *State) save(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	switchStates := make([]netswitch.NetSwitch, 0, len(s.netSwitches))
	for _, ns := range s.netSwitches {
		switchStates = append(switchStates, ns)
	}
	return enc.Encode(switchStates)
}

func (s *State) Save() {
	log.Printf("Saving state...")
	if err := s.save(s.filename); err != nil {
		log.Printf("failed saving state: %v", err)
	}
}

func (s *State) dispatch(cmd Command) (err error) {
	switch cmd.CommandType {
	case CMD_STATE_SYNC:
		wg := sync.WaitGroup{}
		var errs error
		for _, ns := range s.netSwitches {
			if mid := cmd.MachineId; mid != nil && *mid != ns.MachineId {
				continue
			}
			wg.Add(1)
			go func(ns netswitch.NetSwitch) {
				if err := ns.Sync(); err != nil {
					if errs == nil {
						errs = err
					} else {
						errs = fmt.Errorf("%v; %v", errs, err)
					}
				}
				wg.Done()
			}(ns)
		}
		wg.Wait()
		if errs != nil {
			return fmt.Errorf("state watch: %v", errs)
		}
		break
	default:
		log.Fatalf("unknown cmd '%v', terminating...", cmd.CommandType)
	}
	return
}

func (s *State) DispatchLoop() {
	for {
		select {
		case cmd := <-ch:
			cmd.Error <- s.dispatch(cmd)
		}
	}
}

func PingLoop() {
	for {
		select {
		case <-time.After(global.STATE_SYNC_PERIOD):
			cmd := Command{
				CommandType: CMD_STATE_SYNC,
				Error:       make(chan error),
			}
			ch <- cmd
			err := <-cmd.Error
			if err != nil {
				log.Printf("state watch err: %v", err)
			}
		}
	}
}

func Init(user, key, stateFile string) (err error) {

	client := &http.Client{}
	if client.Jar, err = cookiejar.New(nil); err != nil {
		return
	}
	if err := Login(client, user, key); err != nil {
		return fmt.Errorf("login: %v", err)
	}
	if state, err = LoadState(stateFile, client); err != nil {
		return fmt.Errorf("load state: %v", err)
	}

	log.Printf("netswitches: %v", state.netSwitches)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("received signal %v", sig)
			// sig is a ^C, handle it
			state.Save()
			os.Exit(1)
		}
	}()

	go state.DispatchLoop()
	go PingLoop()

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
		ns := state.netSwitches[id]
		ns.On = cmdStr == CMD_ON
		state.netSwitches[id] = ns
		cmd := Command{
			CommandType: CMD_STATE_SYNC,
			MachineId:   &id,
			Error:       make(chan error),
		}
		ch <- cmd
		if err := <-cmd.Error; err != nil {
			return fmt.Errorf("cmd dispatch: %v", err)
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
	apiUrl = *flag.String("apiUrl", "http://localhost:8080/api", "Url of the fabsmith api (http or https)")
	user := flag.String("id", "user", "id")
	key := flag.String("key", "user", "key")
	stateFile := flag.String("stateFile", "state.json", "switches are stateful but they loose state on reset")
	flag.Parse()
	if err := Init(*user, *key, *stateFile); err != nil {
		log.Fatalf("Init: %v", err)
	}

	http.HandleFunc("/machines/", RunCommand)

	if err := http.ListenAndServe(":7070", nil); err != nil {
		state.Save()
	}
}
