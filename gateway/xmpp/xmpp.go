package xmpp

import (
	"crypto/tls"
	"fmt"
	"github.com/FabLabBerlin/localmachines/gateway/global"
	"github.com/mattn/go-xmpp"
	"log"
	"strings"
)

type Xmpp struct {
	ch   chan string
	talk *xmpp.Client
}

func NewXmpp(server, user, pass string) (x *Xmpp, err error) {
	x = &Xmpp{
		ch: make(chan string, 10),
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
				x.ch <- v.Text
			case xmpp.Presence:
				if global.XMPP_DEBUG {
					fmt.Println(v.From, v.Show)
				}
			}
		}
	}()
}

func (x *Xmpp) Recv() <-chan string {
	return x.ch
}

func (x *Xmpp) Send(remote, text string) (err error) {
	_, err = x.talk.Send(xmpp.Chat{
		Remote: remote,
		Type:   "chat",
		Text:   text,
	})
	return
}
