package user_permissions

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const TABLE_NAME = "permission"

type Permission struct {
	Id        int64 `orm:"pk"`
	UserId    int64
	MachineId int64
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

func Create(userId, machineId int64) (err error) {
	permission := Permission{
		UserId:    userId,
		MachineId: machineId,
	}

	o := orm.NewOrm()
	_, _, err = o.ReadOrCreate(&permission, "UserId", "MachineId")

	return
}

func Delete(userId, machineId int64) (err error) {
	p := Permission{
		UserId:    userId,
		MachineId: machineId,
	}

	o := orm.NewOrm()

	if err = o.Read(&p, "UserId", "MachineId"); err != nil {
		return
	}

	_, err = o.Delete(&p)
	return
}

func Update(userId int64, permissions *[]Permission) error {

	// Delete all existing permissions of the user
	o := orm.NewOrm()
	beego.Info("Attempting to delete user permissions row...")
	_, err := o.QueryTable(TABLE_NAME).
		Filter("UserId", userId).Delete()
	if err != nil {
		return fmt.Errorf("Error deleting: %v", err)
	}

	// If there are no permissions, do nothing
	if len(*permissions) <= 0 {
		return nil
	}

	// Create new permissions
	if _, err = o.InsertMulti(len(*permissions), permissions); err != nil {
		return fmt.Errorf("Failed to insert permissions: %v", err)
	}

	return nil
}
