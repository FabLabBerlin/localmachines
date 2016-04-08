package global

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/satori/go.uuid"
	"sync"
	"time"
)

var (
	mu sync.Mutex
	// responses are matched here to the RPC requests.  We don't want to have
	// much blocking here, therefore the channels are buffered (capacity 1) and
	// all reads/writes must happen asynchronously.
	responses            map[string]chan xmpp.Message
	xmppClient *xmpp.Xmpp
	xmppServerConfigured bool
)

func XmppClient() *xmpp.Xmpp {
	return xmppClient
}

func XmppServerConfigured() bool {
	return xmppServerConfigured
}

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

func SendXmppCommand(location *locations.Location, command string, machineId int64) (err error) {
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
