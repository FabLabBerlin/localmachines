package machine

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/easylab-lib/xmpp"
	"github.com/FabLabBerlin/easylab-lib/xmpp/commands"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

const (
	NETSWITCH_TYPE_MFI = "mfi"
	// Empty value as unspecified type/custom
	NETSWITCH_TYPE_CUSTOM = ""
)

var (
	ApplyToBilling       func(update redis.MachinesUpdate)
	xmppServerConfigured bool
	xmppClient           *xmpp.Xmpp
)

func init() {
	server := beego.AppConfig.String("XmppServer")
	xmppServerConfigured = server != ""
	if xmppServerConfigured {
		user := beego.AppConfig.String("XmppUser")
		pass := beego.AppConfig.String("XmppPass")
		debugPrint := func(f string, argv ...interface{}) {
			beego.Info(fmt.Sprintf(f, argv...))
		}
		xmppClient = xmpp.NewXmpp(server, user, pass, debugPrint)
		xmppClient.Run()
		go func() {
			for {
				select {
				case msg := <-xmppClient.Recv():
					if err := xmppDispatch(msg); err != nil {
						beego.Info("xmpp dispatch:", err)
					}
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
	case commands.GATEWAY_ALLOWS_USERS_FROM_IP:
		if msg.Data.IpAddress != "" {
			if err := locations.SetLocalIp(
				msg.Data.LocationId,
				msg.Data.IpAddress,
			); err != nil {
				beego.Error("FetchLocalIpsTask: location=", msg.Data.LocationId, ": could not save local ip: ", err)
			}
		} else {
			beego.Error("FetchLocalIpsTask: location=", msg.Data.LocationId, ": empty ip")
		}
	case commands.GATEWAY_APPLIED_CONFIG_0, commands.GATEWAY_APPLIED_CONFIG_1, commands.GATEWAY_APPLIED_CONFIG_2:
		update := redis.MachinesUpdate{
			LocationId: msg.Data.LocationId,
			MachineId:  msg.Data.MachineId,
			UserId:     msg.Data.UserId,
			Command:    msg.Data.Command,
		}
		if msg.Data.Error {
			update.Error = msg.Data.ErrorMessage
		} else {
			switch msg.Data.Command {
			case commands.GATEWAY_APPLIED_CONFIG_0:
				update.Info = "Connected with Gateway ✔"
			case commands.GATEWAY_APPLIED_CONFIG_1:
				update.Info = "Wifi is configured on switch ✔"
			case commands.GATEWAY_APPLIED_CONFIG_2:
				update.Info = "Switch is reconfigured and rebooting now..."
			}
		}
		if err := redis.PublishMachinesUpdate(update); err != nil {
			beego.Error("publish machines update:", err)
		}
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
	case commands.GATEWAY_SENDS_NETSWITCH_STATUS:
		beego.Info("GATEWAY_SENDS_NETSWITCH_STATUS:", msg.Data)
		if err := SetNetswitchLastPing(msg.Data.MachineId, time.Now()); err != nil {
			beego.Error("SetNetswitchLastPing(", msg.Data.MachineId, "):", err)
		}
	case commands.GATEWAY_SUCCESS_ON:
		update := redis.MachinesUpdate{
			LocationId: msg.Data.LocationId,
			MachineId:  msg.Data.MachineId,
			UserId:     msg.Data.UserId,
			Info:       "Successfully turned on machine",
			Command:    commands.GATEWAY_SUCCESS_ON,
		}
		ApplyToBilling(update)
		if err := redis.PublishMachinesUpdate(update); err != nil {
			beego.Error("publish machines update:", err)
		}
	case commands.GATEWAY_SUCCESS_OFF:
		update := redis.MachinesUpdate{
			LocationId: msg.Data.LocationId,
			MachineId:  msg.Data.MachineId,
			UserId:     msg.Data.UserId,
			Info:       "Successfully turned off machine",
			Command:    commands.GATEWAY_SUCCESS_OFF,
		}
		ApplyToBilling(update)
		if err := redis.PublishMachinesUpdate(update); err != nil {
			beego.Error("publish machines update:", err)
		}
	case commands.GATEWAY_FAIL_ON, commands.GATEWAY_FAIL_OFF:
		update := redis.MachinesUpdate{
			LocationId: msg.Data.LocationId,
			MachineId:  msg.Data.MachineId,
			UserId:     msg.Data.UserId,
			Command:    msg.Data.Command,
		}
		if msg.Data.Command == commands.GATEWAY_FAIL_ON {
			update.Error = "Failed to turn on machine"
		} else {
			update.Error = "Failed to turn off machine"
		}
		if strings.Contains(msg.Data.ErrorMessage, "unreachable") {
			update.Error += " (host unreachable)"
		}
		if strings.Contains(msg.Data.ErrorMessage, "no route to host") {
			update.Error += " (no route to host)"
		}
		if msg.Data.Command == commands.GATEWAY_FAIL_OFF {
			ApplyToBilling(update)
		}
		if err := redis.PublishMachinesUpdate(update); err != nil {
			beego.Error("publish machines update:", err)
		}
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

func (this *Machine) On(userId int64) error {
	return this.turn(ON, userId)
}

func (this *Machine) Off(userId int64) error {
	return this.turn(OFF, userId)
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

func SetNetswitchLastPing(id int64, last time.Time) (err error) {
	query := `
	UPDATE machines
	SET netswitch_last_ping = ?
	WHERE id = ?
	`
	o := orm.NewOrm()
	_, err = o.Raw(query, last.UTC().Format("2006-01-02 15:04:05"), id).Exec()
	return
}

func (this *Machine) turn(onOrOff ON_OR_OFF, userId int64) (err error) {
	beego.Info("Attempt to turn NetSwitch ", onOrOff, ", machine ID", this.Id,
		", NetswitchHost: ", this.NetswitchHost)

	if this.NetswitchConfigured() {
		if xmppClient != nil {
			return this.turnXmpp(onOrOff, userId)
		} else {
			return fmt.Errorf("xmpp client is nil!")
		}
	} else {
		update := redis.MachinesUpdate{
			LocationId: this.LocationId,
			MachineId:  this.Id,
			UserId:     userId,
		}
		if onOrOff == ON {
			update.Info = "Successfully turned on machine"
			update.Command = commands.GATEWAY_SUCCESS_ON
		} else {
			update.Info = "Successfully turned off machine"
			update.Command = commands.GATEWAY_SUCCESS_OFF
		}
		ApplyToBilling(update)
		err = redis.PublishMachinesUpdate(update)
		if err != nil {
			return fmt.Errorf("publish machines update: %v", err)
		}
	}
	return
}

func (this *Machine) turnXmpp(onOrOff ON_OR_OFF, userId int64) (err error) {
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
			UserId:     userId,
		},
	})
	return
}

func (this *Machine) NetswitchApplyConfig(userId int64) (err error) {
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
			UserId:     userId,
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

	// We return always nil.  If things fail, we log them.
	return nil
}
