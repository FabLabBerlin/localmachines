package netswitch

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/satori/go.uuid"
	"net/http"
	"regexp"
	"strings"
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

type Mapping struct {
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
	orm.RegisterModel(new(Mapping))
}

func (m *Mapping) TableName() string {
	return TABLE_NAME
}

func GetAllMappingsAt(locationId int64) (mappings []*Mapping, err error) {
	o := orm.NewOrm()

	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return nil, fmt.Errorf("new query builder: %v", err)
	}

	qb.Select(TABLE_NAME + ".*").
		From(TABLE_NAME).
		InnerJoin("machines").
		On(TABLE_NAME + ".machine_id = machines.id").
		Where("machines.location_id = ?")

	_, err = o.Raw(qb.String(), locationId).
		QueryRows(&mappings)
	return
}

func CreateMapping(machineId int64) (int64, error) {
	mapping := Mapping{}
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

	if err = xmppReinit(); err != nil {
		return 0, fmt.Errorf("xmpp reinit: %v", err)
	}
	return id, nil
}

func GetMapping(machineId int64) (*Mapping, error) {
	mapping := Mapping{}
	mapping.MachineId = machineId
	o := orm.NewOrm()
	err := o.Read(&mapping, "MachineId")
	return &mapping, err
}

func DeleteMapping(machineId int64) (err error) {
	mapping, err := GetMapping(machineId)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	if _, err = o.Delete(mapping); err != nil {
		return fmt.Errorf("orm delete: %v", err)
	}

	if err = xmppReinit(); err != nil {
		return fmt.Errorf("xmpp reinit: %v", err)
	}
	return
}

func (this *Mapping) Update() (err error) {
	o := orm.NewOrm()

	// Check for duplicate host entries
	this.Host = strings.TrimSpace(this.Host)
	if this.Host != "" {
		netswitch := Mapping{}
		query := fmt.Sprintf("SELECT * FROM %v WHERE host=? AND sensor_port=? AND id<>?",
			netswitch.TableName())
		var nsms []Mapping
		num, err := o.Raw(query, this.Host, this.SensorPort, this.Id).QueryRows(&nsms)
		if err != nil {
			return fmt.Errorf("failed to query db: %v", err)
		}
		if num > 0 {
			return fmt.Errorf("Found %v machines with same netswitch host", num)
		}
	}

	if _, err = o.Update(this); err != nil {
		return fmt.Errorf("update: %v", err)
	}
	if err = xmppReinit(); err != nil {
		return fmt.Errorf("xmpp reinit: %v", err)
	}
	return
}

func xmppReinit() (err error) {
	if xmppServerConfigured {
		if err = sendXmppCommand("reinit", 0); err != nil {
			return fmt.Errorf("send xmpp cmd: %v", err)
		}
	}
	return
}

func (this *Mapping) On() error {
	return this.turn(ON)
}

func (this *Mapping) Off() error {
	return this.turn(OFF)
}

func (this *Mapping) turn(onOrOff ON_OR_OFF) (err error) {
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

func (this *Mapping) turnHttp(onOrOff ON_OR_OFF) (err error) {
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

func (this *Mapping) turnXmpp(onOrOff ON_OR_OFF) (err error) {
	return sendXmppCommand(string(onOrOff), this.MachineId)
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
