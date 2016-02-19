package metrics

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type Realtime struct {
	LocationId     int64
	Activations    int64
	ActiveMachines int
	Users          int64
}

func NewRealtime(locationId int64) (rt *Realtime, err error) {
	rt = &Realtime{
		LocationId: locationId,
	}
	if rt.Activations, err = rt.activations(); err != nil {
		return nil, fmt.Errorf("activations count: %v", err)
	}
	if rt.ActiveMachines, err = rt.activeMachines(); err != nil {
		return nil, fmt.Errorf("active machines count: %v", err)
	}
	if rt.Users, err = rt.users(); err != nil {
		return nil, fmt.Errorf("users count: %v", err)
	}
	return
}

// activations number of users that have no membership or paid membership.
func (rt *Realtime) activations() (n int64, err error) {
	query := `
		SELECT Count(*) AS N
		FROM   purchases
		       LEFT JOIN user_membership
		         ON user_membership.user_id = purchases.user_id
		       LEFT JOIN membership
		         ON user_membership.membership_id = membership.id
		WHERE  type = 'activation'
		       AND (monthly_price > 0 OR monthly_price IS NULL)
		       AND purchases.location_id = ?
	`
	var maps []orm.Params
	o := orm.NewOrm()
	if _, err = o.Raw(query, rt.LocationId).Values(&maps); err != nil {
		return
	}
	return strconv.ParseInt(maps[0]["N"].(string), 10, 64)
}

func (rt *Realtime) activeMachines() (n int, err error) {
	all, err := purchases.GetActiveActivations()
	if err != nil {
		return
	}
	for _, a := range all {
		if a.Purchase.LocationId == rt.LocationId {
			n++
		}
	}
	return
}

// users number that have no membership or paid membership.
func (rt *Realtime) users() (n int64, err error) {
	query := `
		SELECT Count(*) AS N
		FROM   user
		       LEFT JOIN user_membership
		         ON user_membership.user_id = user.id
		       LEFT JOIN membership
		         ON user_membership.membership_id = membership.id
		       JOIN user_locations
		         ON user.id = user_locations.user_id
		WHERE  (monthly_price > 0 OR monthly_price IS NULL)
		       AND user_locations.location_id = ?
	`
	var maps []orm.Params
	o := orm.NewOrm()
	if _, err = o.Raw(query, rt.LocationId).Values(&maps); err != nil {
		return
	}
	return strconv.ParseInt(maps[0]["N"].(string), 10, 64)
}
