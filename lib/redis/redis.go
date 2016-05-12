package redis

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

const (
	IdleTimeoutSeconds     = 240
	machinesUpdateChPrefix = "machines_updates"
	NumberMaxIdle          = 10
)

var (
	pool *redis.Pool
)

func init() {
	dsn := beego.AppConfig.String("SessionProviderConfig")
	pool = &redis.Pool{
		MaxIdle:     NumberMaxIdle,
		IdleTimeout: IdleTimeoutSeconds * time.Second,
		Dial: func() (c redis.Conn, err error) {
			return redis.Dial("tcp", dsn)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func GetPool() *redis.Pool {
	return pool
}

func GetPubSubConn() (psc redis.PubSubConn) {
	psc = redis.PubSubConn{
		Conn: pool.Get(),
	}
	return
}

func MachinesUpdateCh(locationId int64) string {
	return machinesUpdateChPrefix + "-" + strconv.FormatInt(locationId, 10)
}

func PublishMachinesUpdate(locationId int64) (err error) {
	beego.Info("PublishMachinesUpdate()")
	ch := MachinesUpdateCh(locationId)
	_, err = pool.Get().Do("PUBLISH", ch, "")
	return
}
