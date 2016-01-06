package global

import "time"

// Actually timeout and period could be set dynamically depending on the current latencies
const (
	MAX_SYNC_RETRIES   = 3
	STATE_SYNC_TIMEOUT = 9 * time.Second
	STATE_SYNC_PERIOD  = 3 * STATE_SYNC_TIMEOUT
)

var (
	ApiUrl        string
	StateFilename string
)

func init() {
	if STATE_SYNC_TIMEOUT.Seconds() >= STATE_SYNC_PERIOD.Seconds() {
		panic("timeout should be smaller than the watch period")
	}
}
