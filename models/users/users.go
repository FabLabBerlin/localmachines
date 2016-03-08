package users

import (
	"errors"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"regexp"
	"strings"
	"time"
)

// Regular expression for email spec : RFC 5322
var _EXP_EMAIL = regexp.MustCompile(`(?i)[-a-z0-9~!$%^*_=+}{\'?]+(\.[-a-z0-9~!$%^*_=+}{\'?]+)*@([a-z0-9_][-a-z0-9_]*(\.[-a-z0-9_]+)*\.(aero|arpa|biz|com|coop|edu|gov|info|int|mil|museum|name|net|org|pro|travel|mobi|[a-z][a-z])|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,5})?`)

type User struct {
	Id          int64  `orm:"auto";"pk"`
	FirstName   string `orm:"size(100)"`
	LastName    string `orm:"size(100)"`
	Username    string `orm:"size(100)"`
	Email       string `orm:"size(100)"`
	InvoiceAddr string `orm:"type(text)"`
	ShipAddr    string `orm:"type(text)"`
	// ClientId is the Fastbill User Id
	ClientId    int64
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

func init() {
	orm.RegisterModel(new(User))
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
	if err := user.CheckEmail(); err != nil {
		return 0, err
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

			user.Id = id
			user.Created = time.Now()
			_, err = o.Update(user)
			if err != nil {
				beego.Error("Failed to update user create time:", err)
				return 0, fmt.Errorf("Failed to update user create time: %v", err)
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

func (user *User) CheckEmail() (err error) {
	return checkEmail(user.Email)
}

func checkEmail(email string) (err error) {
	if strings.TrimSpace(email) == "" {
		return errors.New("Email field should not be blank")
	}
	if !_EXP_EMAIL.MatchString(email) {
		return errors.New("Invalid email")
	}
	return
}

func (user *User) GetRole() user_roles.Role {
	return user_roles.Role(user.UserRole)
}

// Update user
func (user *User) Update() error {
	o := orm.NewOrm()

	// Check if not last admin in case UserRole is not admin
	if user.GetRole() != user_roles.ADMIN {
		numAdmins, err := o.QueryTable(user.TableName()).
			Filter("UserRole", user_roles.ADMIN).Count()
		if err != nil {
			return err
		}
		beego.Trace("Number of admins:", numAdmins)
		if numAdmins <= 1 {
			// Check if the user we are updating is that 1 last admin
			userIsAdmin := o.QueryTable(user.TableName()).
				Filter("id", user.Id).
				Filter("user_role", user_roles.ADMIN).Exist()
			if userIsAdmin {
				return errors.New("Only one admin left")
			}
		}
	}

	// Validate email
	if err := user.CheckEmail(); err != nil {
		return err
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

// Loads user data from database into User struct
func GetUser(userId int64) (*User, error) {
	user := User{Id: userId}
	o := orm.NewOrm()
	err := o.Read(&user)
	if err != nil {
		return nil, fmt.Errorf("Could not get user data, ID:", userId)
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