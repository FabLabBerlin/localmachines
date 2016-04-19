package request_response

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/xmpp"
	"github.com/satori/go.uuid"
	"log"
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

type DispatchFunc func(msg xmpp.Message) (ipAddress string, err error)

func NewDispatcher(server, user, pass string, dispatch DispatchFunc) (d *Dispatcher) {
	log.Printf("NewDispatcher(%v, %v, ...)", server, user)
	d = &Dispatcher{
		xmppClient: xmpp.NewXmpp(server, user, pass),
		responses:  make(map[string]chan xmpp.Message),
	}
	d.xmppClient.Run()
	go func() {
		for {
			select {
			case msg := <-d.xmppClient.Recv():
				if msg.Data.IsRequest {
					log.Printf("incoming request of type %v", msg.Data.Command)
					if dispatch != nil {
						ipAddress, err := dispatch(msg)
						if err != nil {
							log.Printf("xmpp dispatcher: dispatch: %v", err)
						}
						log.Printf("dispatch(msg) called")
						response := xmpp.Message{
							Remote:    msg.Remote,
							Data:      msg.Data,
						}
						response.Data.IsRequest = false
						response.Data.IpAddress = ipAddress
						response.Data.Error = err != nil
						if err := d.xmppClient.Send(response); err != nil {
							log.Printf("xmpp: failed to send response")
						}
					} else {
						log.Printf("dispatcher: got request but no dispatch function")
					}
				} else {
					resp := msg
					d.mu.Lock()
					tid := resp.Data.TrackingId
					select {
					case d.responses[tid] <- resp:
					default:
						log.Printf("package already received: tid: %v", tid)
					}
					d.mu.Unlock()
				}
				break
			}
		}
	}()
	return
}

func (d *Dispatcher) SendXmppCommand(locationId int64, xmppId, command string, machineId int64, payload string) (ipAddress string, err error) {
	trackingId := uuid.NewV4().String()
	d.mu.Lock()
	d.responses[trackingId] = make(chan xmpp.Message, 1)
	respCh := d.responses[trackingId]
	d.mu.Unlock()
	log.Printf("calling xmpp client send")
	err = d.xmppClient.Send(xmpp.Message{
		Remote: xmppId,
		Data: xmpp.Data{
			IsRequest:  true,
			Command:    command,
			MachineId:  machineId,
			TrackingId: trackingId,
			LocationId: locationId,
			Payload:    payload,
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
