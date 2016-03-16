package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Permission struct {
	Id        int64 `orm:"pk"`
	UserId    int64
	MachineId int64
}

func (this *Permission) TableName() string {
	return "permission"
}

func init() {
	orm.RegisterModel(new(Permission))
}

// Returns which machines user is enabled to use
func GetUserPermissions(userId int64) (*[]Permission, error) {
	var permissions []Permission
	o := orm.NewOrm()
	_, err := o.QueryTable("permission").
		Filter("user_id", userId).All(&permissions)
	if err != nil {
		return nil, err
	}
	return &permissions, nil
}

func CreateUserPermission(userId, machineId int64) error {
	beego.Trace("Creating user permission")
	permission := Permission{}
	permission.UserId = userId
	permission.MachineId = machineId
	beego.Trace(permission)

	o := orm.NewOrm()
	created, id, err := o.ReadOrCreate(&permission, "UserId", "MachineId")
	if err != nil {
		return err
	}

	if created {
		beego.Info("User permission created:", id)
	} else {
		beego.Warning("User permission already exists:", id)
	}

	return nil
}

func DeleteUserPermission(userId, machineId int64) error {
	p := Permission{}
	p.UserId = userId
	p.MachineId = machineId

	var err error

	o := orm.NewOrm()
	err = o.Read(&p, "UserId", "MachineId")
	if err != nil {
		return err
	}

	var num int64

	num, err = o.Delete(&p)
	if err != nil {
		return err
	}

	beego.Trace("Num permissions deleted:", num)
	return nil
}

func UpdateUserPermissions(userId int64, permissions *[]Permission) error {

	// Delete all existing permissions of the user
	p := Permission{}
	o := orm.NewOrm()
	beego.Info("Attempting to delete user permissions row...")
	num, err := o.QueryTable(p.TableName()).
		Filter("UserId", userId).Delete()
	if err != nil {
		beego.Error("Error:", err)
		//return err
	}
	beego.Trace("Deleted num permissions:", num)

	// If there are no permissions, do nothing
	if len(*permissions) <= 0 {
		return nil
	}

	// Create new permissions
	num, err = o.InsertMulti(len(*permissions), permissions)
	if err != nil {
		beego.Error("Failed to insert permissions")
		return err
	}
	beego.Trace("Inserted num permissions:", num)

	return nil
}
