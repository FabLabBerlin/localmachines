package xmpp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/mattn/go-xmpp"
	"log"
	"strings"
)

type Xmpp struct {
	ch   chan Message
	talk *xmpp.Client
}

type Message struct {
	Remote string
	Data   Data
}

type Data struct {
	Command    string
	MachineId  int64
	TrackingId int64
	Error      bool
}

func NewXmpp(server, user, pass string) (x *Xmpp, err error) {
	x = &Xmpp{
		ch: make(chan Message, 10),
	}

	xmpp.DefaultConfig = tls.Config{
		ServerName:         serverName(server),
		InsecureSkipVerify: global.XMPP_TLS_INSECURE_SKIP_VERIFY,
	}

	options := xmpp.Options{
		Host:          server,
		User:          user,
		Password:      pass,
		NoTLS:         global.XMPP_NO_TLS,
		Debug:         global.XMPP_DEBUG,
		Session:       global.XMPP_USE_SERVER_SESSION,
		Status:        STATUS,
		StatusMessage: STATUS_MESSAGE,
	}

	x.talk, err = options.NewClient()

	return
}

func serverName(host string) string {
	return strings.Split(host, ":")[0]
}

const (
	STATUS         = "xa"
	STATUS_MESSAGE = "I for one welcome our new codebot overlords."
)

func (x *Xmpp) Run() {
	go func() {
		for {
			chat, err := x.talk.Recv()
			if err != nil {
				log.Fatal(err)
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				fmt.Println(v.Remote, v.Text)

				var data Data
				err := json.Unmarshal([]byte(v.Text), &data)
				if err != nil {
					log.Printf("xmpp: %v", err)
					log.Printf("remote was: '%v'", v.Remote)
					log.Printf("text was: '%v'", v.Text)
				} else {
					x.ch <- Message{
						Remote: v.Remote,
						Data:   data,
					}
				}
			case xmpp.Presence:
				if global.XMPP_DEBUG {
					fmt.Println(v.From, v.Show)
				}
			}
		}
	}()
}

func (x *Xmpp) Recv() <-chan Message {
	return x.ch
}

func (x *Xmpp) Send(msg Message) (err error) {
	buf, err := json.Marshal(msg.Data)
	if err != nil {
		return
	}
	_, err = x.talk.Send(xmpp.Chat{
		Remote: msg.Remote,
		Type:   "chat",
		Text:   string(buf),
	})
	return
}
