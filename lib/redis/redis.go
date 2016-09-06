// packages redis manages a connection pool.
package redis

import (
	"encoding/json"
	"fmt"
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
				PublishMachinesUpdate(MachinesUpdate{
					LocationId: loc.Id,
				})
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

type MachinesUpdate struct {
	LocationId int64
	MachineId  int64
	UserId     int64
	Error      string
	Info       string
	Warning    string
}

func PublishMachinesUpdate(update MachinesUpdate) (err error) {
	beego.Info("PublishMachinesUpdate()")
	ch := MachinesUpdateCh(update.LocationId)
	conn := pool.Get()
	defer conn.Close()
	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("json marshal: %v", err)
	}
	_, err = conn.Do("PUBLISH", ch, string(data))
	return
}
