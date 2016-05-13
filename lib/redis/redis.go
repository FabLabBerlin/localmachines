package redis

import (
	"github.com/FabLabBerlin/localmachines/models/locations"
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
	go func() {
		for {
			beego.Info("pub update for all locations")
			<-time.After(time.Minute)
			locs, err := locations.GetAll()
			if err != nil {
				beego.Error("locations get all:", err)
				continue
			}
			for _, loc := range locs {
				PublishMachinesUpdate(loc.Id)
			}
		}
	}()
}

func GetPool() *redis.Pool {
	return pool
}

func GetPoolConn() (psc redis.Conn) {
	return pool.Get()
}

func MachinesUpdateCh(locationId int64) string {
	return machinesUpdateChPrefix + "-" + strconv.FormatInt(locationId, 10)
}

func PublishMachinesUpdate(locationId int64) (err error) {
	beego.Info("PublishMachinesUpdate()")
	ch := MachinesUpdateCh(locationId)
	conn := pool.Get()
	defer conn.Close()
	_, err = conn.Do("PUBLISH", ch, "")
	return
}
