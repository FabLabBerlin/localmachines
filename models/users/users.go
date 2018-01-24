package users

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/FabLabBerlin/localmachines/models/user_locations"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// Regular expression for email spec : RFC 5322
var _EXP_EMAIL = regexp.MustCompile(`(?i)[-a-z0-9~!$%^*_=+}{\'?]+(\.[-a-z0-9~!$%^*_=+}{\'?]+)*@([a-z0-9_][-a-z0-9_]*(\.[-a-z0-9_]+)*\.(aero|arpa|biz|com|coop|edu|gov|info|int|mil|museum|name|net|org|pro|travel|mobi|[a-z][a-z])|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,5})?`)

type User struct {
	Id          int64
	FirstName   string `orm:"size(100)"`
	LastName    string `orm:"size(100)"`
	Username    string `orm:"size(100)"`
	Email       string `orm:"size(100)"`
	InvoiceAddr string `orm:"type(text)"`
	ShipAddr    string `orm:"type(text)"`
	// ClientId is the Fastbill User Id
	ClientId           int64
	B2b                bool
	Company            string `orm:"size(100)"`
	VatUserId          string `orm:"size(100)"`
	VatRate            int
	SuperAdmin         bool
	Created            time.Time `orm:"type(datetime)"`
	Comments           string    `orm:"type(text)"`
	Phone              string    `orm:"size(50)"`
	ZipCode            string    `orm:"size(100)"`
	City               string    `orm:"size(100)"`
	CountryCode        string    `orm:"size(2)"`
	NoAutoInvoicing    bool
	FastbillTemplateId int64
	EuDelivery         bool
	// University fields
	StudentId        *string
	SecurityBriefing *string
}

func (this *User) TableName() string {
	return "user"
}

func init() {
	orm.RegisterModel(new(User))
}

var ErrEmailExists = errors.New("User with the same email exists")
var ErrUsernameExists = errors.New("User with the same username exists")

// Attempt to create user, do not complain if it already exists
func CreateUser(user *User) (userId int64, er error) {

	var err error
	var num int64
	var query string
	type UserIdStruct struct {
		Id int64
	}

	if user.SuperAdmin {
		return 0, errors.New("attempted to create user with super admin rights")
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
			return 0, ErrUsernameExists
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
			return 0, ErrEmailExists
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

// Update user
func (user *User) Update() error {
	o := orm.NewOrm()

	// Validate email
	if err := user.CheckEmail(); err != nil {
		return err
	}

	if len(strings.TrimSpace(user.FirstName)) == 0 ||
		len(strings.TrimSpace(user.LastName)) == 0 {
		return fmt.Errorf("First/last name mustn't be empty")
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
	beego.Trace("Found num users with the same email: %v", num)
	if num > 0 {
		return fmt.Errorf("User with the same email exists: %v", err)
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
		return nil, fmt.Errorf("Could not get user data, ID: %v", userId)
	} else {
		return &user, nil
	}
}

func GetUserByClientId(clientId int64) (*User, error) {
	var user User

	err := orm.NewOrm().
		QueryTable("user").
		Filter("client_id", clientId).
		One(&user)

	return &user, err
}

// Returns an array with all users in the system
func GetAllUsersAt(locationId int64) (users []*User, err error) {
	var all []*User

	o := orm.NewOrm()
	if _, err = o.QueryTable("user").All(&all); err != nil {
		return users, fmt.Errorf("Failed to get all users: %v", err)
	}

	uls, err := user_locations.GetAllForLocation(locationId)
	if err != nil {
		return nil, fmt.Errorf("get all for location: %v", err)
	}
	userIds := make(map[int64]struct{})
	for _, ul := range uls {
		userIds[ul.UserId] = struct{}{}
	}

	users = make([]*User, 0, len(userIds))
	for _, u := range all {
		if _, ok := userIds[u.Id]; ok {
			users = append(users, u)
		}
	}
	return users, nil
}
