package machine

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/FabLabBerlin/localmachines/lib/xmpp/commands"
	"github.com/FabLabBerlin/localmachines/lib/xmpp/request_response"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
	"strings"
)

const (
	NETSWITCH_TYPE_MFI = "mfi"
	// Empty value as unspecified type/custom
	NETSWITCH_TYPE_CUSTOM = ""
)

var (
	xmppServerConfigured bool
	dispatcher *request_response.Dispatcher
)

func init() {
	server := beego.AppConfig.String("XmppServer")
	xmppServerConfigured = server != ""
	if xmppServerConfigured {
		user := beego.AppConfig.String("XmppUser")
		pass := beego.AppConfig.String("XmppPass")
		dispatcher = request_response.NewDispatcher(server, user, pass, dispatch)
	}
}

type ON_OR_OFF string

const (
	ON  ON_OR_OFF = "on"
	OFF ON_OR_OFF = "off"
)

func xmppReinit(location *locations.Location) (err error) {
	if xmppServerConfigured {
		machines, err := GetAllAt(location.Id)
		if err != nil {
			return fmt.Errorf("get machines: %v", err)
		}
		buf, err := json.Marshal(machines)
		if err != nil {
			return fmt.Errorf("json marshal: %v", err)
		}
		if _, err = dispatcher.SendXmppCommand(location.Id, location.XmppId, "reinit", 0, string(buf)); err != nil {
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
		if dispatcher != nil {
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
	_, err = dispatcher.SendXmppCommand(location.Id, location.XmppId, string(onOrOff), this.Id, "")
	return
}

func (this *Machine) NetswitchApplyConfig() (err error) {
	location, err := locations.Get(this.LocationId)
	if err != nil {
		return fmt.Errorf("get location %v: %v", this.LocationId, err)
	}
	_, err = dispatcher.SendXmppCommand(location.Id, location.XmppId, commands.APPLY_CONFIG, this.Id, "")
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
		ipAddress, err := dispatcher.SendXmppCommand(l.Id, l.XmppId, commands.FETCH_LOCAL_IP, 0, "")
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

func dispatch(msg xmpp.Message) (ipAddress string, err error) {
	if cmd := msg.Data.Command; cmd == commands.REQUEST_REINIT {
		locId := msg.Data.LocationId
		l, err := locations.Get(locId)
		if err != nil {
			return "", fmt.Errorf("error getting location %v", locId)
		}
		if err = xmppReinit(l); err != nil {
			return "", fmt.Errorf("xmpp reinit: %v", err)
		}
	} else {
		return "", fmt.Errorf("unknown command %v", cmd)
	}
	return "", nil
}
