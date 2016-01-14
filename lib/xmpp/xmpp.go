package xmpp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/mattn/go-xmpp"
	"log"
	"strings"
	"time"
)

const (
	DEBUG                    = false
	NO_TLS                   = false
	TLS_INSECURE_SKIP_VERIFY = false
	USE_SERVER_SESSION       = true
)

type Xmpp struct {
	ch     chan Message
	talk   *xmpp.Client
	user   string
	server string
}

type Message struct {
	Remote string
	Data   Data
}

type Data struct {
	Command    string
	MachineId  int64
	TrackingId string
	Error      bool
}

func NewXmpp(server, user, pass string) (x *Xmpp, err error) {
	x = &Xmpp{
		ch:     make(chan Message, 10),
		user:   user,
		server: server,
	}

	xmpp.DefaultConfig = tls.Config{
		ServerName:         serverName(server),
		InsecureSkipVerify: TLS_INSECURE_SKIP_VERIFY,
	}

	options := xmpp.Options{
		Host:          server,
		User:          user,
		Password:      pass,
		NoTLS:         NO_TLS,
		Debug:         DEBUG,
		Session:       USE_SERVER_SESSION,
		Status:        STATUS,
		StatusMessage: STATUS_MESSAGE,
	}

	x.talk, err = options.NewClient()

	if err == nil {
		go func() {
			for {
				<-time.After(time.Second)
				log.Printf("Pinnnnggggggggg")
				if err := x.Ping(); err != nil {
					log.Printf("ping errrrrr: %v", err)
				}
			}
		}()
	}

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
				log.Printf("xmpp chat rcvd")
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
				log.Printf("xmpp presence rcvd")
				if DEBUG {
					fmt.Println(v.From, v.Show)
				}
			case xmpp.IQ:
				log.Printf("xmpp iq rcvd: %v", v)
			default:
				log.Printf("1234: unknown msg type")
			}
		}
	}()
}

func (x *Xmpp) Ping() error {
	//tmp := fmt.Sprintf("bla%v", time.Now().Unix())
	//log.Printf("tmp=%v, server=%v\n", tmp, x.server)
	return x.talk.PingC2S(x.user, "api.easylab.io" /*x.server*/)
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
