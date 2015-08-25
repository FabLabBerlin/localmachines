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
	"regexp"
	"strings"
	"time"
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

// Regular expression for email spec : RFC 5322
const _EXP_EMAIL = `(?i)[-a-z0-9~!$%^*_=+}{\'?]+(\.[-a-z0-9~!$%^*_=+}{\'?]+)*@([a-z0-9_][-a-z0-9_]*(\.[-a-z0-9_]+)*\.(aero|arpa|biz|com|coop|edu|gov|info|int|mil|museum|name|net|org|pro|travel|mobi|[a-z][a-z])|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,5})?`

type User struct {
	Id          int64  `orm:"auto";"pk"`
	FirstName   string `orm:"size(100)"`
	LastName    string `orm:"size(100)"`
	Username    string `orm:"size(100)"`
	Email       string `orm:"size(100)"`
	InvoiceAddr string `orm:"type(text)"`
	ShipAddr    string `orm:"type(text)"`
	// ClientId is the Fastbill User Id
	ClientId    int
	B2b         bool
	Company     string `orm:"size(100)"`
	VatUserId   string `orm:"size(100)"`
	VatRate     int
	UserRole    string    `orm:"size(100)"`
	Created     time.Time `orm:"type(datetime)"`
	Comments    string    `orm:"type(text)"`
	Phone       string    `orm:"size(50)"`
	ZipCode     string    `orm:"size(100)"`
	City        string    `orm:"size(100)"`
	CountryCode string    `orm:"size(2)"`
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

func (this *Auth) TableName() string {
	return "auth"
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

// Attempt to create user, do not complain if it already exists
func CreateUser(user *User) (userId int64, er error) {

	var err error
	var num int64
	var query string
	type UserIdStruct struct {
		Id int64
	}

	o := orm.NewOrm()

	// Validate email
	if user.Email == "" {
		return 0, errors.New("Email field should not be blank")
	}
	exp, err := regexp.Compile(_EXP_EMAIL)
	if err != nil {
		return 0, fmt.Errorf("Failed to compile rexex: %v", err)
	}
	if !exp.MatchString(user.Email) {
		return 0, errors.New("Invalid email")
	}

	// Check if user with the same username does exist
	if user.Username != "" {
		query = fmt.Sprintf("SELECT id FROM %s WHERE username=?", user.TableName())
		var userIds []UserIdStruct
		num, err = o.Raw(query, user.Username).QueryRows(&userIds)
		if err != nil {
			return 0, fmt.Errorf("Failed to exec query: %v", err)
		}
		if num > 0 {
			return 0, errors.New("User with the same username exists")
		}
	}

	// Check if user with the same email does exist
	if user.Email != "" {
		query = fmt.Sprintf("SELECT id FROM %s WHERE email=?", user.TableName())
		var userEmailIds []UserIdStruct
		num, err = o.Raw(query, user.Email).QueryRows(&userEmailIds)
		if err != nil {
			return 0, fmt.Errorf("Failed to exec query: %v", err)
		}
		if num > 0 {
			return 0, errors.New("User with the same email exists")
		}
	}

	var created bool
	var id int64
	if created, id, err = o.ReadOrCreate(user, "Email"); err == nil {
		if created {
			beego.Info("Created user with ID", id)

			// Update created time.
			// Still using a hack here as Beego timezone support is foggy.
			updQuery := fmt.Sprintf("UPDATE %s SET created=? WHERE id=?",
				user.TableName())
			_, err = o.Raw(updQuery,
				time.Now().Format("2006-01-02 15:04:05"), id).Exec()
			if err != nil {
				beego.Error("Failed to update user creation time")
			}
		} else {
			beego.Info("User already exists with ID", id,
				"and email", user.Email)
			return 0, fmt.Errorf("User exists")
		}
		return id, nil
	} else {
		return 0, fmt.Errorf("Could not ReadOrCreate user: %v", err)
	}
}

// Authenticate username and password, return user ID on success
func AuthenticateUser(username, password string) (int64, error) {
	authModel := Auth{}
	o := orm.NewOrm()
	var err error
	query := "SELECT hash, salt FROM auth INNER JOIN user ON auth.user_id = user.id"
	if strings.Contains(username, "@") {
		beego.Info("query by email")
		err = o.Raw(query+" WHERE user.email = ?", username).QueryRow(&authModel)
	} else {
		beego.Info("query by user name")
		err = o.Raw(query+" WHERE user.username = ?", username).QueryRow(&authModel)
	}

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
		var err error
		query := o.QueryTable("user")
		if strings.Contains(username, "@") {
			err = query.Filter("email", username).One(&user, "id")
		} else {
			err = query.Filter("username", username).One(&user, "id")
		}
		if err != nil {
			beego.Error("Could not get user ID for user", username)
			return 0, errors.New("Failed to authenticate")
		}
		return user.Id, nil
	} else {
		return 0, errors.New("Failed to authenticate")
	}
}

// Authenticate user by using NFC uid
// TODO: provide some kind of basic crypto
func AuthenticateUserUid(uid string) (string, int64, error) {
	var err error

	// uid can not be empty or less than 4 chars
	if uid == "" || len(uid) < 4 {
		err = errors.New("Invalid NFC UID")
		return "", 0, err
	}

	auth := Auth{}
	auth.NfcKey = uid
	o := orm.NewOrm()

	// Get user ID
	err = o.Read(&auth, "NfcKey")
	if err != nil {
		beego.Error(err)
		return "", 0, errors.New("Failed to read auth table")
	}

	// Get user name
	user := User{}
	user.Id = auth.UserId
	err = o.Read(&user, "Id")
	if err != nil {
		return "", 0, errors.New(
			fmt.Sprintf("Failed to read user table: %v", err))
	}

	return user.Username, user.Id, nil
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

func AuthUpdateNfcUid(userId int64, nfcUid string) error {
	var err error
	var num int64
	o := orm.NewOrm()
	auth := Auth{UserId: userId}
	err = o.Read(&auth)
	authRecordMissing := err == orm.ErrNoRows
	if err != nil && !authRecordMissing {
		return fmt.Errorf("Missing auth record: %v", err)
	}

	// No update required if the UIDs already match
	if auth.NfcKey == nfcUid {
		beego.Warning("This UID is already assigned to the user")
		return nil
	}

	// Check if another user uses the same UID
	num, err = o.QueryTable(auth.TableName()).Filter("NfcKey", nfcUid).Count()
	if err != nil {
		beego.Warning("Failed to get matching auth records")
	}
	if num > 0 {
		return errors.New("Auth records with the UID exist")
	}

	auth.NfcKey = nfcUid
	num, err = o.Update(&auth, "NfcKey")
	if err != nil {
		return fmt.Errorf("Failed to update: %v", err)
	}
	beego.Trace("Update affected num rows:", num)
	return nil
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
			// Check if the user we are updating is that 1 last admin
			userIsAdmin := o.QueryTable(user.TableName()).
				Filter("id", user.Id).
				Filter("user_role", "admin").Exist()
			if userIsAdmin {
				return errors.New("Only one admin left")
			}
		}
	}

	// Validate email
	if user.Email == "" {
		return errors.New("Email field should not be blank")
	}
	exp, err := regexp.Compile(_EXP_EMAIL)
	if err != nil {
		return fmt.Errorf("Failed to compile rexex: %v", err)
	}
	if !exp.MatchString(user.Email) {
		return errors.New("Invalid email")
	}

	// Check for duplicate username and email entries
	var query string
	query = fmt.Sprintf("SELECT id FROM %s WHERE username=? AND id!=?",
		user.TableName())
	type UserIdStruct struct {
		Id int64
	}
	var sameUsernameUserIds []UserIdStruct
	num, err := o.Raw(query, user.Username, user.Id).
		QueryRows(&sameUsernameUserIds)
	if err != nil {
		return fmt.Errorf("Failed to query db: %v", err)
	}
	beego.Trace("Found num users with the same username:", num)
	if num > 0 {
		return fmt.Errorf("User with the same username exists: %v", err)
	}
	query = fmt.Sprintf("SELECT id FROM %s WHERE email=? AND id!=?",
		user.TableName())
	var sameEmailUserIds []UserIdStruct
	num, err = o.Raw(query, user.Email, user.Id).
		QueryRows(&sameEmailUserIds)
	if err != nil {
		return fmt.Errorf("Failed to query db: %v", err)
	}
	beego.Trace("Found num users with the same email:", num)
	if num > 0 {
		return fmt.Errorf("User with the same email exists", err)
	}

	// TODO: Check email regex

	// Update the user finally
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

	// Delete all user machine permissions associated with this user
	perm := Permission{}
	num, err = o.QueryTable(perm.TableName()).Filter("user_id",
		userId).Delete()
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to delete user machine permissions: %v", err))
	}
	beego.Trace("Deleted num user machine permissions:", num)

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
