package machine

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/FabLabBerlin/localmachines/lib/xmpp/commands"
	"github.com/FabLabBerlin/localmachines/models/locations"
)

func TaskPingNetswitches() (err error) {
	var locId int64 = 1

	loc, err := locations.Get(locId)
	if err != nil {
		return fmt.Errorf("get location %v: %v", locId, err)
	}

	return pingNetswitchesAt(loc)
}

func pingNetswitchesAt(loc *locations.Location) (err error) {
	return xmppClient.Send(xmpp.Message{
		Remote: loc.XmppId,
		Data: xmpp.Data{
			Command:    commands.FETCH_NETSWITCH_STATUS,
			LocationId: loc.Id,
		},
	})
}
