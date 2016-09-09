/*
global configuration variables.
*/
package global

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Actually timeout and period could be set dynamically depending on the current latencies
const (
	MAX_SYNC_RETRIES   = 3
	STATE_SYNC_TIMEOUT = 9 * time.Second
	STATE_SYNC_PERIOD  = 3 * STATE_SYNC_TIMEOUT
)

var (
	Cfg            Config
	ServerPrefix   *string
	ServerJabberId string
)

type Config struct {
	Main struct {
		LocationId int64
	}
	XMPP struct {
		Server string
		User   string
		Pass   string
	}
}

func init() {
	if STATE_SYNC_TIMEOUT.Seconds() >= STATE_SYNC_PERIOD.Seconds() {
		panic("timeout should be smaller than the watch period")
	}
}

func DebugHttp(msg string) {
	if ServerPrefix == nil {
		return
	}
	l := strconv.FormatInt(Cfg.Main.LocationId, 10)
	resp, err := http.PostForm(*ServerPrefix+"/api/locations/debug?location="+l,
		url.Values{"message": {msg}})
	if err != nil {
		log.Printf("DebugHttp: %v", err)
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
}
