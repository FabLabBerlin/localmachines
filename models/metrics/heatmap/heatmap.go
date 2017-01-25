package heatmap

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/countries"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"net/http"
	"strings"
	"time"
)

// HEATMAP: http://sunng87.github.io/heatcanvas/openstreetmap.html
// geo coding eph 2 redis: http://wiki.openstreetmap.org/wiki/Nominatim

// TaskGeoCodeUsers can be run hourly, it'll geo code only the things that
// aren't done yet. Redis cache ~ inf
func TaskGeoCodeUsers() (err error) {
	beego.Info("TaskGeoCodeUsers()")
	ls, err := locations.GetAll()
	if err != nil {
		return
	}
	for _, l := range ls {
		us, err := users.GetAllUsersAt(l.Id)
		if err != nil {
			return err
		}
		i := 0
		for _, u := range us {
			k := fmt.Sprintf("geocode(%v)", u.Id)
			if redis.Exists(k) {
				continue
			}
			var tmp interface{}
			err := redis.Cached(k, 2592000, tmp, geoCode(*u))
			if err != nil {
				beego.Error("redis.Cached:", err)
			}
			<-time.After(10 * time.Second)
			i++
			if i > 10 {
				break
			}
		}
	}
	return
}

type Coordinate struct {
	Lat float64 `json:"lat,string"`
	Lon float64 `json:"lon,string"`
}

func geoCode(u users.User) func() (coord interface{}, err error) {
	return func() (coord interface{}, err error) {
		url := "https://nominatim.openstreetmap.org/search/"
		url += "?format=json"
		country, ok := countries.GetByCode(u.CountryCode)
		if ok {
			url += "&countrycodes=" + strings.ToLower(country.Code)
		} else {
			url += "&countrycodes=de"
		}
		url += "&q=" + fmt.Sprintf("%v, %v %v", u.InvoiceAddr, u.ZipCode, u.City)
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)
		var res []Coordinate
		if err = dec.Decode(&res); err != nil {
			return
		}
		if len(res) == 0 {
			return coord, nil
		}
		beego.Info("res[0]=", res[0])
		return res[0], nil
	}
}
