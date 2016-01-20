package models

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/satori/go.uuid"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var (
	mu sync.Mutex
	// responses are matched here to the RPC requests.  We don't want to have
	// much blocking here, therefore the channels are buffered (capacity 1) and
	// all reads/writes must happen asynchronously.
	responses   map[string]chan xmpp.Message
	xmppClient  *xmpp.Xmpp
	xmppGateway string
)

func init() {
	if server := beego.AppConfig.String("XmppServer"); server != "" {
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

type NetSwitchMapping struct {
	Id        int64 `orm:"auto";"pk"`
	MachineId int64
	UrlOn     string `orm:"size(255)"`
	UrlOff    string `orm:"size(255)"`
	// Host and Sensor Port are only relevant for mfi switches
	Host       string `orm:"size(255)"`
	SensorPort int
	Xmpp       bool
}

func init() {
	orm.RegisterModel(new(NetSwitchMapping))
}

func (m *NetSwitchMapping) TableName() string {
	return "netswitch"
}

func GetAllNetSwitchMapping() (ms []*NetSwitchMapping, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(NetSwitchMapping).TableName()).All(&ms)
	return
}

func CreateNetSwitchMapping(machineId int64) (int64, error) {
	mapping := NetSwitchMapping{}
	mapping.MachineId = machineId
	mapping.SensorPort = 1
	o := orm.NewOrm()
	created, id, err := o.ReadOrCreate(&mapping, "MachineId")
	if err != nil {
		return 0, err
	}
	if created {
		beego.Info("Created new NetSwitch mapping, machine ID:", id)
	} else {
		beego.Warning("A NetSwitch mapping exists, machine ID:", id)
	}
	return id, nil
}

func GetNetSwitchMapping(machineId int64) (*NetSwitchMapping, error) {
	mapping := NetSwitchMapping{}
	mapping.MachineId = machineId
	o := orm.NewOrm()
	err := o.Read(&mapping, "MachineId")
	if err != nil {
		return &mapping, err
	}
	return &mapping, nil
}

func DeleteNetSwitchMapping(machineId int64) error {
	mapping, err := GetNetSwitchMapping(machineId)
	if err != nil {
		return err
	}

	var num int64

	o := orm.NewOrm()
	num, err = o.Delete(mapping)
	if err != nil {
		return err
	}

	beego.Trace("Affected num rows:", num)
	return nil
}

func (this *NetSwitchMapping) Update() (err error) {
	o := orm.NewOrm()
	if _, err = o.Update(this); err != nil {
		return fmt.Errorf("update: %v", err)
	}
	if err = this.sendXmppCommand("reinit", 0); err != nil {
		return fmt.Errorf("send xmpp cmd: %v", err)
	}
	return
}

func (this *NetSwitchMapping) On() error {
	return this.turn(ON)
}

func (this *NetSwitchMapping) Off() error {
	return this.turn(OFF)
}

func (this *NetSwitchMapping) turn(onOrOff ON_OR_OFF) (err error) {
	beego.Info("Attempt to turn NetSwitch ", onOrOff, ", machine ID", this.MachineId)
	if this.Xmpp {
		if xmppClient != nil {
			return this.turnXmpp(onOrOff)
		} else {
			return fmt.Errorf("xmpp client is nil!")
		}
	} else {
		return this.turnHttp(onOrOff)
	}
}

func (this *NetSwitchMapping) turnHttp(onOrOff ON_OR_OFF) (err error) {
	var resp *http.Response

	if onOrOff == ON {
		resp, err = http.Get(this.UrlOn)
	} else {
		resp, err = http.Get(this.UrlOff)
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

func (this *NetSwitchMapping) turnXmpp(onOrOff ON_OR_OFF) (err error) {
	return this.sendXmppCommand(string(onOrOff), this.MachineId)
}

func (this *NetSwitchMapping) sendXmppCommand(command string, machineId int64) (err error) {
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
