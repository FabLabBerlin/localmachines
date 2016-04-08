package machine

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/global"
	"github.com/astaxie/beego"
	"strings"
)

const (
	NETSWITCH_TYPE_MFI = "mfi"
	// Empty value as unspecified type/custom
	NETSWITCH_TYPE_CUSTOM = ""
)

type ON_OR_OFF string

const (
	ON  ON_OR_OFF = "on"
	OFF ON_OR_OFF = "off"
)

func xmppReinit(location *locations.Location) (err error) {
	if global.XmppServerConfigured() {
		if err = global.SendXmppCommand(location, "reinit", 0); err != nil {
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
		", NetswitchHost: ", this.NetswitchHost)

	mfiConfigured := this.NetswitchType == NETSWITCH_TYPE_MFI &&
		len(this.NetswitchHost) > 0

	customConfigured := this.NetswitchType == NETSWITCH_TYPE_CUSTOM &&
		len(strings.TrimSpace(this.NetswitchUrlOn)) > 0 &&
		len(strings.TrimSpace(this.NetswitchUrlOff)) > 0
	if mfiConfigured || customConfigured {
		if global.XmppClient() != nil {
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
	return global.SendXmppCommand(location, string(onOrOff), this.Id)
}

func (this *Machine) NetswitchApplyConfig() (err error) {
	location, err := locations.Get(this.LocationId)
	if err != nil {
		return fmt.Errorf("get location %v: %v", this.LocationId, err)
	}
	return global.SendXmppCommand(location, "apply_config", this.Id)
}
