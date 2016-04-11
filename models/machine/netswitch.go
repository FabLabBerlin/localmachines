package machine

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/FabLabBerlin/localmachines/lib/xmpp/commands"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
	"strings"
	"sync"
	"time"
)

const (
	NETSWITCH_TYPE_MFI = "mfi"
	// Empty value as unspecified type/custom
	NETSWITCH_TYPE_CUSTOM = ""
)

var (
	mu sync.Mutex
	// responses are matched here to the RPC requests.  We don't want to have
	// much blocking here, therefore the channels are buffered (capacity 1) and
	// all reads/writes must happen asynchronously.
	responses            map[string]chan xmpp.Message
	xmppClient           *xmpp.Xmpp
	xmppServerConfigured bool
)

func init() {
	server := beego.AppConfig.String("XmppServer")
	xmppServerConfigured = server != ""
	if xmppServerConfigured {
		user := beego.AppConfig.String("XmppUser")
		pass := beego.AppConfig.String("XmppPass")
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

func xmppReinit(location *locations.Location) (err error) {
	if xmppServerConfigured {
		if _, err = sendXmppCommand(location, "reinit", 0); err != nil {
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

func (this *Machine) NetswitchConfigured() bool {
	return this.NetswitchCustomConfigured() || this.NetswitchMfiConfigured()
}

func (this *Machine) NetswitchCustomConfigured() bool {
	return this.NetswitchType == NETSWITCH_TYPE_CUSTOM &&
		len(strings.TrimSpace(this.NetswitchUrlOn)) > 0 &&
		len(strings.TrimSpace(this.NetswitchUrlOff)) > 0
}

func (this *Machine) NetswitchMfiConfigured() bool {
	return this.NetswitchType == NETSWITCH_TYPE_MFI &&
		len(this.NetswitchHost) > 0
}

func (this *Machine) turn(onOrOff ON_OR_OFF) (err error) {
	beego.Info("Attempt to turn NetSwitch ", onOrOff, ", machine ID", this.Id,
		", NetswitchHost: ", this.NetswitchHost)

	if this.NetswitchConfigured() {
		if xmppClient != nil {
			return this.turnXmpp(onOrOff)
		} else {
			return fmt.Errorf("xmpp client is nil!")
		}
	}
	return
}

func (this *Machine) turnXmpp(onOrOff ON_OR_OFF) (err error) {
	location, err := locations.Get(this.LocationId)
	if err != nil {
		return fmt.Errorf("get location %v: %v", this.LocationId, err)
	}
	_, err = sendXmppCommand(location, string(onOrOff), this.Id)
	return
}

func (this *Machine) NetswitchApplyConfig() (err error) {
	location, err := locations.Get(this.LocationId)
	if err != nil {
		return fmt.Errorf("get location %v: %v", this.LocationId, err)
	}
	_, err = sendXmppCommand(location, commands.APPLY_CONFIG, this.Id)
	return
}

func sendXmppCommand(location *locations.Location, command string, machineId int64) (ipAddress string, err error) {
	trackingId := uuid.NewV4().String()
	mu.Lock()
	responses[trackingId] = make(chan xmpp.Message, 1)
	respCh := responses[trackingId]
	mu.Unlock()
	err = xmppClient.Send(xmpp.Message{
		Remote: location.XmppId,
		Data: xmpp.Data{
			Command:    command,
			MachineId:  machineId,
			TrackingId: trackingId,
			LocationId: location.Id,
		},
	})
	if err != nil {
		return "", fmt.Errorf("send: %v", err)
	}
	select {
	case resp := <-respCh:
		if resp.Data.Error {
			err = fmt.Errorf("some error occurred")
		} else {
			ipAddress = resp.Data.IpAddress
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

func FetchLocalIpsTask() error {
	beego.Info("Running FetchLocalIpsTask")

	locs, err := locations.GetAll()
	if err != nil {
		return fmt.Errorf("get locations: %v", err)
	}

	for _, l := range locs {
		if l.XmppId == "" {
			continue
		}
		ipAddress, err := sendXmppCommand(l, commands.FETCH_LOCAL_IP, 0)
		if err != nil {
			beego.Error("FetchLocalIpsTask: location=", l.Id, ":", err)
		}
		if ipAddress != "" {
			l.LocalIp = ipAddress
			if err := locations.SetLocalIp(l.Id, ipAddress); err != nil {
				beego.Error("FetchLocalIpsTask: location=", l.Id, ": could not save local ip: ", err)
			}
		} else {
			beego.Error("FetchLocalIpsTask: location=", l.Id, ": empty ip")
		}
	}

	// We return always nil.  If things fail, we log them.
	return nil
}
