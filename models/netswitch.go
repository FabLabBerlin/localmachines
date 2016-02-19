package models

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
	"net/http"
	"regexp"
	"sync"
	"time"
)

const TABLE_NAME = "netswitch"

var (
	mu sync.Mutex
	// responses are matched here to the RPC requests.  We don't want to have
	// much blocking here, therefore the channels are buffered (capacity 1) and
	// all reads/writes must happen asynchronously.
	responses            map[string]chan xmpp.Message
	xmppClient           *xmpp.Xmpp
	xmppGateway          string
	xmppServerConfigured bool
)

func init() {
	server := beego.AppConfig.String("XmppServer")
	xmppServerConfigured = server != ""
	if xmppServerConfigured {
		user := beego.AppConfig.String("XmppUser")
		pass := beego.AppConfig.String("XmppPass")
		xmppGateway = beego.AppConfig.String("XmppGateway")
		xmppClient = xmpp.NewXmpp(server, user, pass)
		xmppClient.Run()

		responses = make(map[string]chan xmpp.Message)
		go func() {
			for {
				select {
				case resp := <-xmppClient.Recv():
					mu.Lock()
					tid := resp.Data.TrackingId
					select {
					case responses[tid] <- resp:
					default:
						beego.Error("package already received: tid:", tid)
					}
					mu.Unlock()
					break
				}
			}
		}()
	}
}

type ON_OR_OFF string

const (
	ON  ON_OR_OFF = "on"
	OFF ON_OR_OFF = "off"
)

func xmppReinit() (err error) {
	if xmppServerConfigured {
		if err = sendXmppCommand("reinit", 0); err != nil {
			return fmt.Errorf("send xmpp cmd: %v", err)
		}
	}
	return
}

func (this *Machine) On() error {
	return this.turn(ON)
}

func (this *Machine) Off() error {
	return this.turn(OFF)
}

func (this *Machine) turn(onOrOff ON_OR_OFF) (err error) {
	beego.Info("Attempt to turn NetSwitch ", onOrOff, ", machine ID", this.Id,
		", NetswitchXmpp: ", this.NetswitchXmpp, ", NetswitchHost: ", this.NetswitchHost)
	beego.Info("whole obj: %v", *this)
	if this.NetswitchXmpp {
		if xmppClient != nil {
			return this.turnXmpp(onOrOff)
		} else {
			return fmt.Errorf("xmpp client is nil!")
		}
	} else {
		return this.turnHttp(onOrOff)
	}
}

func (this *Machine) turnHttp(onOrOff ON_OR_OFF) (err error) {
	var resp *http.Response

	if onOrOff == ON {
		resp, err = http.Get(this.NetswitchUrlOn)
	} else {
		resp, err = http.Get(this.NetswitchUrlOff)
	}

	if err != nil {
		// Work around custom HTTP status code the switch returns: "AhmaSwitch"
		matched, _ := regexp.MatchString("malformed HTTP status code", err.Error())
		if !matched {
			return fmt.Errorf("Failed to send NetSwitch %v request: %v", onOrOff, err)
		}
	}

	beego.Trace(resp)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			beego.Error("Bad Status Code:", resp.StatusCode)
			return errors.New("Bad Status Code")
		}
	}

	return nil
}

func (this *Machine) turnXmpp(onOrOff ON_OR_OFF) (err error) {
	return sendXmppCommand(string(onOrOff), this.Id)
}

func sendXmppCommand(command string, machineId int64) (err error) {
	trackingId := uuid.NewV4().String()
	mu.Lock()
	responses[trackingId] = make(chan xmpp.Message, 1)
	respCh := responses[trackingId]
	mu.Unlock()
	err = xmppClient.Send(xmpp.Message{
		Remote: xmppGateway,
		Data: xmpp.Data{
			Command:    command,
			MachineId:  machineId,
			TrackingId: trackingId,
		},
	})
	if err != nil {
		return fmt.Errorf("send: %v", err)
	}
	select {
	case resp := <-respCh:
		if resp.Data.Error {
			err = fmt.Errorf("some error occurred")
		} else {
			err = nil
		}
		break
	case <-time.After(20 * time.Second):
		err = fmt.Errorf("timeout")
	}

	mu.Lock()
	delete(responses, trackingId)
	mu.Unlock()

	return
}
