/*
user_locations package provides location based user permissions.
*/
package user_locations

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/astaxie/beego/orm"
)

const (
	TABLE_NAME = "user_locations"
)

type UserLocation struct {
	Id         int64
	LocationId int64
	Location   *locations.Location `orm:"-"`
	UserId     int64
	UserRole   string
	Archived   bool
}

func (ul *UserLocation) TableName() string {
	return TABLE_NAME
}

func Create(ul *UserLocation) (id int64, err error) {
	return orm.NewOrm().Insert(ul)
}

func GetAllForLocation(locationId int64) (uls []*UserLocation, err error) {
	_, err = orm.NewOrm().QueryTable(TABLE_NAME).
		Filter("location_id", locationId).
		Exclude("archived", 1).
		All(&uls)
	return
}

func GetAllForUser(userId int64) (uls []*UserLocation, err error) {
	ls, err := locations.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all user locations: %v", err)
	}
	_, err = orm.NewOrm().QueryTable(TABLE_NAME).
		Filter("user_id", userId).
		Exclude("archived", 1).
		All(&uls)
	for _, ul := range uls {
		for _, l := range ls {
			if ul.LocationId == l.Id {
				ul.Location = l
				break
			}
		}
	}
	return
}

func (ul *UserLocation) Update() (err error) {
	_, err = orm.NewOrm().Update(ul)
	return
}

func init() {
	orm.RegisterModel(new(UserLocation))
}
