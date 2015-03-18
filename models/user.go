package models

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/scrypt"
	"io"
	"errors"
)

// cf. http://stackoverflow.com/a/23039768/485185
const (
    PW_SALT_BYTES = 32
    PW_HASH_BYTES = 64
)

type User struct {
	Id          int    `orm:"auto";"pk"`
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
}

type Auth struct {
	UserId   int    `orm:"auto"`
	NfcKey   string `orm:"size(100)"`
	Hash     string `orm:"size(300)"`
	Salt     string `orm:"size(100)"`
}

type UserRoles struct {
	UserId int  `orm:"auto"`
	Admin  bool `orm:"size(100)"`
	Staff  bool `orm:"size(100)"`
	Member bool `orm:"size(100)"`
}

func init() {
	orm.RegisterModel(new(User), new(Auth), new(UserRoles))
}

// Authenticate username and password, return user ID on success
func AuthenticateUser(username, password string) (int, error) {
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

// Loads user data from database into User struct
func GetUser(userId int) (*User, error) {
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
	beego.Trace("Got num users: ", num)
	return users, nil
}

// Loads user roles from database into UserRoles struct
func GetUserRoles(userId int) (*UserRoles, error) {
	userRoles := UserRoles{UserId: userId}
	o := orm.NewOrm()
	err := o.Read(&userRoles)
	if err != nil {
		beego.Error("Could not get user roles, ID:", userId)
		return nil, err
	} else {
		return &userRoles, nil
	}
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

