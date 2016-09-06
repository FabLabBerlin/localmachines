package machine

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/FabLabBerlin/localmachines/lib/xmpp/commands"
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
	xmppClient           *xmpp.Xmpp
)

func init() {
	fmt.Printf("netswitch.go#init\n")
	beego.Info("netswitch.go#init")
	server := beego.AppConfig.String("XmppServer")
	xmppServerConfigured = server != ""
	if xmppServerConfigured {
		beego.Info("Initializing XMPP Client...")
		fmt.Printf("initializing xmpp client...\n")
		user := beego.AppConfig.String("XmppUser")
		pass := beego.AppConfig.String("XmppPass")
		xmppClient = xmpp.NewXmpp(server, user, pass)
		xmppClient.Run()
		go func() {
			for {
				fmt.Printf("bla123543543534\n")
				select {
				case msg := <-xmppClient.Recv():
					if err := xmppDispatch(msg); err != nil {
						beego.Info("xmpp dispatch:", err)
					}
					beego.Info("INCOMING PKG!!!!:", msg)
					fmt.Printf("INC PKG!!\n")
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

func xmppDispatch(msg xmpp.Message) (err error) {
	switch msg.Data.Command {
	case commands.GATEWAY_REQUESTS_MACHINE_LIST:
		ms, err := GetAllAt(msg.Data.LocationId)
		if err != nil {
			return fmt.Errorf("machines.GetAllAt: %v", err)
		}

		raw, err := json.Marshal(ms)
		if err != nil {
			return fmt.Errorf("json marshal: %v", err)
		}

		location, err := locations.Get(msg.Data.LocationId)
		if err != nil {
			return fmt.Errorf("get location: %v", err)
		}

		if err = xmppClient.Send(xmpp.Message{
			Remote: location.XmppId,
			Data: xmpp.Data{
				Command:    commands.SERVER_SENDS_MACHINE_LIST,
				LocationId: msg.Data.LocationId,
				Raw:        string(raw),
			},
		}); err != nil {
			return fmt.Errorf("SERVER_SENDS_MACHINE_LIST: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown command '%v'", err)
	}
	return
}

func xmppReinit(location *locations.Location) (err error) {
	if xmppServerConfigured {
		if err = xmppClient.Send(xmpp.Message{
			Remote: location.XmppId,
			Data: xmpp.Data{
				Command:    commands.REINIT,
				LocationId: location.Id,
			},
		}); err != nil {
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
	err = xmppClient.Send(xmpp.Message{
		Remote: location.XmppId,
		Data: xmpp.Data{
			Command:    string(onOrOff),
			LocationId: location.Id,
			MachineId:  this.Id,
		},
	})
	return
}

func (this *Machine) NetswitchApplyConfig() (err error) {
	location, err := locations.Get(this.LocationId)
	if err != nil {
		return fmt.Errorf("get location %v: %v", this.LocationId, err)
	}
	err = xmppClient.Send(xmpp.Message{
		Remote: location.XmppId,
		Data: xmpp.Data{
			Command:    commands.APPLY_CONFIG,
			LocationId: location.Id,
			MachineId:  this.Id,
		},
	})
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
		if err = xmppClient.Send(xmpp.Message{
			Remote: l.XmppId,
			Data: xmpp.Data{
				Command:    commands.FETCH_LOCAL_IP,
				LocationId: l.Id,
			},
		}); err != nil {
			beego.Error("FetchLocalIpsTask: location=", l.Id, ":", err)
		}
	}
	/*ipAddress, err := dispatcher.SendXmppCommand(l, commands.FETCH_LOCAL_IP, 0)
		if err != nil {
			beego.Error("FetchLocalIpsTask: location=", l.Id, ":", err)
		}
		beego.Info("received ip", ipAddress, "for location", l.Id)
		if ipAddress != "" {
			l.LocalIp = ipAddress
			if err := locations.SetLocalIp(l.Id, ipAddress); err != nil {
				beego.Error("FetchLocalIpsTask: location=", l.Id, ": could not save local ip: ", err)
			}
		} else {
			beego.Error("FetchLocalIpsTask: location=", l.Id, ": empty ip")
		}&/
	}*/

	// We return always nil.  If things fail, we log them.
	return nil
}
