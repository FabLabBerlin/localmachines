package redis

import (
	"encoding/json"
	"github.com/astaxie/beego"
)

type Seconds int

func Cached(key string, expire Seconds, o interface{}, f ...func() interface{}) {
	var buf []byte

	c := GetPoolConn()
	defer c.Close()

	c.Send("GET", key)
	c.Flush()

	v, err := c.Receive()
	if err != nil {
		beego.Info("redis cached: receive:", err)
		goto uncached
	}
	if v != nil {
		buf = v.([]byte)

		if err := json.Unmarshal(buf, &o); err != nil {
			beego.Info("redis cached: unmarshal:", err)
		}

		return
	}

uncached:
	if len(f) == 1 {
		o = f[0]()
		buf, err := json.Marshal(o)
		if err != nil {
			beego.Info("redis cached: marshal:", err)
			return
		}
		if err := c.Send("SET", key, buf, "EX", expire); err != nil {
			beego.Info("redis cached: set:", err)
		}
		c.Flush()
	}
	return
}
