package metrics

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type Realtime struct {
	Activations int64
	Users       int64
}

func NewRealtime() (rt *Realtime, err error) {
	rt = &Realtime{}
	if rt.Activations, err = rt.activations(); err != nil {
		return nil, fmt.Errorf("activations count: %v", err)
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
		       AND monthly_price > 0 OR monthly_price IS NULL
	`
	var maps []orm.Params
	o := orm.NewOrm()
	if _, err = o.Raw(query).Values(&maps); err != nil {
		return
	}
	return strconv.ParseInt(maps[0]["N"].(string), 10, 64)
}

// users number that have no membership or paid membership.
func (rt *Realtime) users() (n int64, err error) {
	query := `
		SELECT Count(*) AS N
		FROM   user
		       LEFT JOIN user_membership
		         ON user_membership.user_id = purchases.user_id
		       LEFT JOIN membership
		         ON user_membership.membership_id = membership.id
		WHERE  monthly_price > 0 OR monthly_price IS NULL
	`
	var maps []orm.Params
	o := orm.NewOrm()
	if _, err = o.Raw(query).Values(&maps); err != nil {
		return
	}
	return strconv.ParseInt(maps[0]["N"].(string), 10, 64)
}
