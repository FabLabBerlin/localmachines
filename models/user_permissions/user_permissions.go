package user_permissions

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const TABLE_NAME = "permission"

type Permission struct {
	Id         int64 `orm:"pk"`
	UserId     int64
	CategoryId int64
}

func (this *Permission) TableName() string {
	return TABLE_NAME
}

func init() {
	orm.RegisterModel(new(Permission))
}

// Returns which machines user is enabled to use
func Get(userId int64) (ps *[]Permission, err error) {
	var permissions []Permission
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).
		Filter("user_id", userId).
		All(&permissions)
	return &permissions, err
}

func Create(userId, categoryId int64) (err error) {
	permission := Permission{
		UserId:     userId,
		CategoryId: categoryId,
	}

	o := orm.NewOrm()
	_, _, err = o.ReadOrCreate(&permission, "UserId", "MachineId")

	return
}

func Delete(userId, categoryId int64) (err error) {
	p := Permission{
		UserId:     userId,
		CategoryId: categoryId,
	}

	o := orm.NewOrm()

	if err = o.Read(&p, "UserId", "MachineId"); err != nil {
		return
	}

	_, err = o.Delete(&p)
	return
}

func Update(userId, locationId int64, permissions *[]Permission) error {

	// Delete all existing permissions of the user
	o := orm.NewOrm()
	beego.Info("Attempting to delete user permissions row...")
	query := `
	DELETE permission
	FROM permission
	JOIN machine_types ON machine_types.id = permission.category_id
	WHERE user_id = ? AND location_id = ?`
	if _, err := o.Raw(query, userId, locationId).Exec(); err != nil {
		return fmt.Errorf("Error deleting: %v", err)
	}

	// If there are no permissions, do nothing
	if len(*permissions) <= 0 {
		return nil
	}

	ids := make([]int64, 0, len(*permissions))
	for _, p := range *permissions {
		ids = append(ids, p.Id)
	}

	ms, err := machine.GetMulti(ids)
	if err != nil {
		return fmt.Errorf("Get machines multi: %v", err)
	}

	for _, m := range ms {
		if m.LocationId != locationId {
			return fmt.Errorf("Wrong location id: %v", err)
		}
	}

	// Create new permissions
	if _, err = o.InsertMulti(len(*permissions), permissions); err != nil {
		return fmt.Errorf("Failed to insert permissions: %v", err)
	}

	return nil
}
