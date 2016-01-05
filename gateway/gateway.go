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
	"sync"
	"time"
)

type CommandType string

const (
	CMD_ON          = "on"
	CMD_OFF         = "off"
	CMD_STATE_WATCH = "state_watch"
)

// Actually timeout and period could be set dynamically depending on the current latencies
const (
	STATE_WATCH_TIMEOUT = time.Second
	STATE_WATCH_PERIOD  = 3 * STATE_WATCH_TIMEOUT
)

func init() {
	if STATE_WATCH_TIMEOUT.Seconds() >= STATE_WATCH_PERIOD.Seconds() {
		panic("timeout should be smaller than the watch period")
	}
}

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

var netSwitches map[int64]NetSwitch

type NetSwitch struct {
	Id        int64
	MachineId int64
	UrlOn     string
	UrlOff    string
	On        bool
}

func (ns NetSwitch) stateWatch() (err error) {
	urlStatus := ns.UrlStatus()
	log.Printf("urlStatus = %v", urlStatus)
	if urlStatus == "" {
		return fmt.Errorf("NetSwitch status url for Machine %v empty", ns.MachineId)
	}
	client := http.Client{
		Timeout: STATE_WATCH_TIMEOUT,
	}
	resp, err := client.Get(urlStatus)
	if err != nil {
		return fmt.Errorf("http get url status: %v", err)
	}
	defer resp.Body.Close()
	mfi := MfiSwitch{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&mfi); err != nil {
		return fmt.Errorf("json decode:", err)
	}
	if mfi.On() != ns.On {
		log.Printf("State for Machine %v must get synchronized", ns.MachineId)
		if ns.On {
			if err = ns.TurnOn(); err != nil {
				return fmt.Errorf("turn on: %v", err)
			}
		} else {
			if err = ns.TurnOff(); err != nil {
				return fmt.Errorf("turn off: %v", err)
			}
		}
	}
	return
}

func (ns NetSwitch) TurnOn() (err error) {
	log.Printf("turn on %v", ns.UrlOn)
	resp, err := http.Get(ns.UrlOn)
	if err != nil {
		return fmt.Errorf("client get: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return
}

func (ns NetSwitch) TurnOff() (err error) {
	log.Printf("turn off %v", ns.UrlOff)
	resp, err := http.Get(ns.UrlOff)
	if err != nil {
		return fmt.Errorf("client get: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return
}

// UrlStatus is a really dirty function.  It's just here for the proof of concept.
func (ns NetSwitch) UrlStatus() string {
	tmp := strings.Split(ns.UrlOn, "//")
	host := strings.Split(tmp[1], "/")[0]
	return "http://" + host + "/sensors/1"
}

//{"sensors":[{"output":1,"power":0.0,"energy":0.0,"enabled":0,"current":0.0,"voltage":233.546874046,"powerfactor":0.0,"relay":1,"lock":0}],"status":"success"}

type MfiSwitch struct {
	Sensors []MfiSensor `json:"sensors"`
	Status  string      `json:"status"`
}

func (swi *MfiSwitch) On() bool {
	relay := swi.Sensors[0].Relay
	switch relay {
	case 0:
		return false
		break
	case 1:
		return true
		break
	}
	log.Fatalf("unknown relay status %v, terminating", relay)
	return false
}

type MfiSensor struct {
	Output      int     `json:"output"`
	Power       float64 `json:"power"`
	Energy      float64 `json:"energy"`
	Enabled     float64 `json:"enabled"`
	Current     float64 `json:"current"`
	Voltage     float64 `json:"voltage"`
	PowerFactor float64 `json:"powerfactor"`
	Relay       int     `json:"relay"`
	Lock        int     `json:"lock"`
}

func Login(client *http.Client, apiUrl, user, key string) (err error) {
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

func Fetch(client *http.Client, apiUrl string) (err error) {
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
	netSwitches = make(map[int64]NetSwitch)
	for _, ns := range nss {
		netSwitches[ns.MachineId] = ns
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
			netSwitches[mid] = netswitch
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
		switchStates = append(switchStates, ns)
	}
	return enc.Encode(switchStates)
}

func SaveState(filename string) {
	log.Printf("Saving state...")
	if err := saveState(filename); err != nil {
		log.Printf("failed saving state: %v", err)
	}
}

func dispatch(cmd Command) (err error) {
	switch cmd.CommandType {
	case CMD_ON:
		ns, ok := netSwitches[*cmd.MachineId]
		if !ok {
			return fmt.Errorf("there's no netswitch for machine id %v",
				cmd.MachineId)
		}
		if err = ns.TurnOn(); err != nil {
			return fmt.Errorf("turn on: %v", err)
		}
		ns.On = true
		netSwitches[*cmd.MachineId] = ns
		break
	case CMD_OFF:
		ns, ok := netSwitches[*cmd.MachineId]
		if !ok {
			return fmt.Errorf("there's no netswitch for machine id %v",
				cmd.MachineId)
		}
		if err = ns.TurnOff(); err != nil {
			return fmt.Errorf("turn off: %v", err)
		}
		ns.On = false
		netSwitches[*cmd.MachineId] = ns
		break
	case CMD_STATE_WATCH:
		wg := sync.WaitGroup{}
		var errs error
		for _, ns := range netSwitches {
			wg.Add(1)
			go func(ns NetSwitch) {
				if err := ns.stateWatch(); err != nil {
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

func DispatchLoop() {
	for {
		select {
		case cmd := <-ch:
			cmd.Error <- dispatch(cmd)
		}
	}
}

func PingLoop() {
	for {
		select {
		case <-time.After(STATE_WATCH_PERIOD):
			cmd := Command{
				CommandType: CMD_STATE_WATCH,
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

func Init(apiUrl, user, key, stateFile string) (err error) {

	client := &http.Client{}
	if client.Jar, err = cookiejar.New(nil); err != nil {
		return
	}
	if err := Login(client, apiUrl, user, key); err != nil {
		return fmt.Errorf("login: %v", err)
	}
	if err := Fetch(client, apiUrl); err != nil {
		return fmt.Errorf("fetch: %v", err)
	}
	if err := loadState(stateFile); err != nil {
		return fmt.Errorf("load state: %v", err)
	}

	log.Printf("netswitches: %v", netSwitches)

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

	go DispatchLoop()
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
		commandType := CommandType(cmdStr)
		cmd := Command{
			CommandType: commandType,
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
