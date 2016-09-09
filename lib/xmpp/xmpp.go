/*
xmpp for easy communication behind firewalls.

This package is both used by the netswitch model and the Gateway server.
*/
package xmpp

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/mattn/go-xmpp"
	"log"
	"math/big"
	"strings"
	"time"
)

const (
	DEBUG                    = false
	NO_TLS                   = false
	TLS_INSECURE_SKIP_VERIFY = false
	USE_SERVER_SESSION       = true
	RECONNECT_INIT_WAIT_TIME = time.Second
)

type Xmpp struct {
	ch         chan Message
	talk       *xmpp.Client
	user       string
	server     string
	options    xmpp.Options
	lastPong   time.Time
	debugPrint func(string, ...interface{})
}

type Message struct {
	Remote string
	Data   Data
}

type Data struct {
	Command             string
	MachineId           int64
	LocationId          int64
	UserId              int64
	IpAddress           string
	NetswitchRelayState string
	NetswitchCurrent    float64
	Raw                 string
	Error               bool
	ErrorMessage        string `json:",omitempty"`
}

func NewXmpp(server, user, pass string, debugPrint func(string, ...interface{})) (x *Xmpp) {
	x = &Xmpp{
		ch:         make(chan Message, 10),
		user:       user,
		server:     server,
		debugPrint: debugPrint,
	}
	xmpp.DefaultConfig = tls.Config{
		ServerName:         serverName(server),
		InsecureSkipVerify: TLS_INSECURE_SKIP_VERIFY,
	}
	x.options = xmpp.Options{
		Host:          server,
		User:          user,
		Password:      pass,
		NoTLS:         NO_TLS,
		Debug:         DEBUG,
		Session:       USE_SERVER_SESSION,
		Status:        STATUS,
		StatusMessage: STATUS_MESSAGE,
	}
	go x.connect()
	return
}

func (x *Xmpp) connect() {
	waitTime := RECONNECT_INIT_WAIT_TIME
	for {
		var err error
		x.debugPrint("Xmpp: connecting to Server...")
		x.talk, err = x.options.NewClient()
		if err == nil {
			x.debugPrint("Xmpp: connected to Server.")
			x.lastPong = time.Now()
			waitTime = RECONNECT_INIT_WAIT_TIME
			for {
				if err := x.Ping(); err != nil {
					x.debugPrint("ping errrrrr: %v", err)
					break
				}
				if time.Now().Sub(x.lastPong) > 20*time.Second {
					x.debugPrint("No Pong since 20 seconds, reconnecting!")
					break
				}
				pingWait(0.3, 3)
			}
			x.talk.Close()
			x.talk = nil
		} else {
			x.debugPrint("Xmpp: error connecting to server: %v", err)
			waitTime *= 2
			if waitTime > 30*time.Second {
				waitTime = 30 * time.Second
			}
		}
		<-time.After(waitTime)
	}
}

// Make sure we don't send to many period patterns that could be used to
// break SSL/TLS using crypto analysis
func pingWait(min, max float64) {
	t := time.Second
	r, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err == nil {
		x := float64(r.Int64()) / 10000
		n := min + x*max
		t = time.Duration(n) * t
	} else {
		log.Printf("endpoints: xmpp: pingWait: %v", err)
	}
	<-time.After(t)
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
			if x.talk == nil {
				<-time.After(time.Second)
				continue
			}
			chat, err := x.talk.Recv()
			if err != nil {
				log.Printf("xmpp: receive: %v", err)
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
					if data.Error {
						log.Printf("error was: '%v'", data.ErrorMessage)
					}
					x.ch <- Message{
						Remote: v.Remote,
						Data:   data,
					}
				}
			case xmpp.Presence:
				if DEBUG {
					fmt.Println(v.From, v.Show)
				}
			case xmpp.IQ:
				if v.Type == "result" {
					x.lastPong = time.Now()
				}
			default:
			}
		}
	}()
}

func (x *Xmpp) Ping() error {
	if x.talk != nil {
		return x.talk.PingC2S(x.user, serverName(x.server))
	} else {
		return fmt.Errorf("xmpp: Ping: xmpp client: not ready")
	}
}

func (x *Xmpp) Recv() <-chan Message {
	return x.ch
}

func (x *Xmpp) Send(msg Message) error {
	if x.talk != nil {
		buf, err := json.Marshal(msg.Data)
		if err != nil {
			return err
		}
		_, err = x.talk.Send(xmpp.Chat{
			Remote: msg.Remote,
			Type:   "chat",
			Text:   string(buf),
		})
		return err
	} else {
		return fmt.Errorf("xmpp: Send: xmpp client: not ready")
	}
}
