package redis

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

type Seconds int

func Cached(key string, expire Seconds, o interface{}, f ...func() (interface{}, error)) error {
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
			return fmt.Errorf("unmarshal:", err)
		}

		return nil
	}

uncached:
	if len(f) == 1 {
		fResult, err := f[0]()
		if err != nil {
			return fmt.Errorf("calling f: %v", err)
		}
		buf, err := json.Marshal(fResult)
		if err != nil {
			return fmt.Errorf("marshal:", err)
		}
		if err := c.Send("SET", key, buf, "EX", expire); err != nil {
			return fmt.Errorf("set:", err)
		}
		c.Flush()
		if err := json.Unmarshal(buf, &o); err != nil {
			return fmt.Errorf("cache miss: unmarshal:", err)
		}
	}
	return nil
}

func Exists(key string) (yes bool) {
	c := GetPoolConn()
	defer c.Close()

	exists, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		beego.Error("exists check:", err)
	}
	c.Flush()

	return exists
}
