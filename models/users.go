package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/scrypt"
	"io"
)

// cf. http://stackoverflow.com/a/23039768/485185
const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

const (
	ADMIN  = "admin"
	STAFF  = "staff"
	MEMBER = "member"
)

type User struct {
	Id          int64  `orm:"auto";"pk"`
	FirstName   string `orm:"size(100)"`
	LastName    string `orm:"size(100)"`
	Username    string `orm:"size(100)"`
	Email       string `orm:"size(100)"`
	InvoiceAddr int
	ShipAddr    int
	ClientId    int
	B2b         bool
	Company     string `orm:"size(100)"`
	VatUserId   string `orm:"size(100)"`
	VatRate     int
	UserRole    string `orm:"size(100)"`
}

func (this *User) TableName() string {
	return "user"
}

type Auth struct {
	UserId int64  `orm:"pk"`
	NfcKey string `orm:"size(100)"`
	Hash   string `orm:"size(300)"`
	Salt   string `orm:"size(100)"`
}

type Permission struct {
	Id        int64 `orm:"pk"`
	UserId    int64
	MachineId int64
}

func (this *Permission) TableName() string {
	return "permission"
}

func init() {
	orm.RegisterModel(new(User), new(Auth), new(Permission))
}

// Authenticate username and password, return user ID on success
func AuthenticateUser(username, password string) (int64, error) {
	authModel := Auth{}
	o := orm.NewOrm()
	err := o.Raw("SELECT hash, salt FROM auth INNER JOIN user ON auth.user_id = user.id WHERE user.username = ?",
		username).QueryRow(&authModel)
	if err != nil {
		beego.Error("Could not read into AuthModel:", err)
		return 0, err
	}
	authModelSalt, err := hex.DecodeString(authModel.Salt)
	if err != nil {
		beego.Error("Could not decode authModel.Salt:", err)
		return 0, err
	}
	hash, err := hash(password, authModelSalt)
	if err != nil {
		beego.Error("Could not calculate hash:")
		return 0, err
	}
	isAuthSuccessful := hex.EncodeToString(hash) == authModel.Hash
	if isAuthSuccessful {
		user := User{}
		err := o.QueryTable("user").Filter("username", username).One(&user, "id")
		if err != nil {
			beego.Error("Could not get user ID for user", username)
			return 0, errors.New("Failed to authenticate")
		}
		return user.Id, nil
	} else {
		return 0, errors.New("Failed to authenticate")
	}
}

func AuthSetPassword(userId int64, password string) error {
	o := orm.NewOrm()
	auth := Auth{UserId: userId}
	err := o.Read(&auth)
	authRecordMissing := err == orm.ErrNoRows
	if err != nil && !authRecordMissing {
		return fmt.Errorf("Read: %v", err)
	}
	salt, err := createSalt()
	if err != nil {
		return fmt.Errorf("createSalt: %v", err)
	}
	auth.UserId = userId
	auth.Salt = hex.EncodeToString(salt)
	hash, err := hash(password, salt)
	if err != nil {
		return fmt.Errorf("createHash: %v", err)
	}
	auth.Hash = hex.EncodeToString(hash)
	if authRecordMissing {
		fmt.Printf("insert: auth: %v", auth)
		_, err = o.Insert(&auth)
	} else {
		fmt.Printf("update: auth: %v", auth)
		_, err = o.Update(&auth)
	}
	return err
}

func DeleteUserAuth(userId int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Auth{UserId: userId})
	return err
}

// Update user
func UpdateUser(user *User) error {
	o := orm.NewOrm()

	// Check if not last admin in case UserRole is not admin
	if user.UserRole != "admin" {
		numAdmins, err := o.QueryTable(user.TableName()).
			Filter("UserRole", "admin").Count()
		if err != nil {
			return err
		}
		beego.Trace("Number of admins:", numAdmins)
		if numAdmins <= 1 {
			return errors.New("Only one admin left")
		}
	}

	if _, err := o.Update(user); err != nil {
		return err
	}

	return nil
}

// Deletes user
func DeleteUser(userId int64) error {
	var num int64
	var err error

	// Delete user
	o := orm.NewOrm()
	num, err = o.Delete(&User{Id: userId})
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete user: %v", err))
	}
	beego.Trace("Deleted num user rows:", num)

	// Delete all user activations along with the user
	act := Activation{}
	num, err = o.QueryTable(act.TableName()).Filter("user_id", userId).Delete()
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete user activations: %v", err))
	}
	beego.Trace("Deleted num user activations:", num)

	// Delete all user memberships along with the user
	umem := UserMembership{}
	num, err = o.QueryTable(umem.TableName()).Filter("user_id", userId).Delete()
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete user memberships: %v", err))
	}
	beego.Trace("Deleted num user memberships:", num)

	return nil
}

// Loads user data from database into User struct
func GetUser(userId int64) (*User, error) {
	user := User{Id: userId}
	o := orm.NewOrm()
	err := o.Read(&user)
	if err != nil {
		beego.Error("Could not get user data, ID:", userId)
		return nil, err
	} else {
		return &user, nil
	}
}

// Returns an array with all users in the system
func GetAllUsers() ([]*User, error) {
	o := orm.NewOrm()
	var users []*User
	num, err := o.QueryTable("user").All(&users)
	if err != nil {
		beego.Error("Failed to get all users")
		return users, errors.New("Failed to get all users")
	}
	beego.Trace("Got num users:", num)
	return users, nil
}

// Returns which machines user is enabled to use
func GetUserPermissions(userId int64) (*[]Permission, error) {
	var permissions []Permission
	o := orm.NewOrm()
	num, err := o.QueryTable("permission").
		Filter("user_id", userId).All(&permissions)
	if err != nil {
		return nil, err
	}
	beego.Trace("Got num permissions:", num)
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
		beego.Warning("User permission already exists:", id)
	} else {
		beego.Info("User permission created:", id)
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
	num, err := o.QueryTable(p.TableName()).
		Filter("UserId", userId).Delete()
	if err != nil {
		return err
	}
	beego.Trace("Deleted num permissions:", num)

	// Create new permissions
	num, err = o.InsertMulti(len(*permissions), permissions)
	if err != nil {
		return err
	}
	beego.Trace("Inserted num permissions:", num)

	return nil
}

// Generate scrypt hash
func hash(password string, salt []byte) ([]byte, error) {
	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, PW_HASH_BYTES)
	if err != nil {
		return []byte{}, err
	}
	return hash, nil
}

// Create salt for scrypt
func createSalt() ([]byte, error) {
	salt := make([]byte, PW_SALT_BYTES)
	_, err := io.ReadFull(rand.Reader, salt)
	return salt, err
}
