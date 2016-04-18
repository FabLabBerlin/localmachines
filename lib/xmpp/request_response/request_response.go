package request_response

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/satori/go.uuid"
	"sync"
	"time"
)

type Dispatcher struct {
	mu sync.Mutex
	// responses are matched here to the RPC requests.  We don't want to have
	// much blocking here, therefore the channels are buffered (capacity 1) and
	// all reads/writes must happen asynchronously.
	responses            map[string]chan xmpp.Message
	xmppClient           *xmpp.Xmpp
}

func NewDispatcher(server, user, pass string) (d *Dispatcher) {
	d = &Dispatcher{
		xmppClient: xmpp.NewXmpp(server, user, pass),
		responses:  make(map[string]chan xmpp.Message),
	}
	d.xmppClient.Run()
	go func() {
		for {
			select {
			case resp := <-d.xmppClient.Recv():
				d.mu.Lock()
				tid := resp.Data.TrackingId
				select {
				case d.responses[tid] <- resp:
				default:
					beego.Error("package already received: tid:", tid)
				}
				d.mu.Unlock()
				break
			}
		}
	}()
	return
}

func (d *Dispatcher) SendXmppCommand(location *locations.Location, command string, machineId int64) (ipAddress string, err error) {
	trackingId := uuid.NewV4().String()
	d.mu.Lock()
	d.responses[trackingId] = make(chan xmpp.Message, 1)
	respCh := d.responses[trackingId]
	d.mu.Unlock()
	err = d.xmppClient.Send(xmpp.Message{
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

	d.mu.Lock()
	delete(d.responses, trackingId)
	d.mu.Unlock()

	return
}
