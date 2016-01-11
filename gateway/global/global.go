package global

import "time"

// Actually timeout and period could be set dynamically depending on the current latencies
const (
	MAX_SYNC_RETRIES              = 3
	STATE_SYNC_TIMEOUT            = 9 * time.Second
	STATE_SYNC_PERIOD             = 3 * STATE_SYNC_TIMEOUT
	XMPP_DEBUG                    = false
	XMPP_NO_TLS                   = false
	XMPP_TLS_INSECURE_SKIP_VERIFY = false
	XMPP_USE_SERVER_SESSION       = false
)

var (
	Cfg Config
)

type Config struct {
	Main struct {
		StateFile string
	}
	API struct {
		Id  string
		Key string
		Url string
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
