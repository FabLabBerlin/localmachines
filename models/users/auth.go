package users

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/scrypt"
	"io"
	"strings"
	"time"
)

// cf. http://stackoverflow.com/a/23039768/485185
const (
	PW_SALT_BYTES      = 32
	PW_HASH_BYTES      = 64
	PW_RESET_KEY_BYTES = 64
)

type Auth struct {
	UserId      int64     `orm:"pk"`
	NfcKey      string    `orm:"size(100)"`
	Hash        string    `orm:"size(300)"`
	Salt        string    `orm:"size(100)"`
	PwResetKey  string    `orm:"size(255)"`
	PwResetTime time.Time `orm:"type(datetime)"`
}

func (this *Auth) TableName() string {
	return "auth"
}

func init() {
	orm.RegisterModel(new(Auth))
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
		return 0, fmt.Errorf("Could not read into AuthModel: %v", err)
	}
	authModelSalt, err := hex.DecodeString(authModel.Salt)
	if err != nil {
		return 0, fmt.Errorf("Could not decode authModel.Salt: %v", err)
	}
	hash, err := hash(password, authModelSalt)
	if err != nil {
		return 0, fmt.Errorf("Could not calculate hash: %v", err)
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
	uid = strings.TrimSpace(uid)

	// uid can not be empty or less than 4 chars
	if uid == "" || len(uid) < 4 {
		return "", 0, errors.New("Invalid NFC UID")
	}

	auth := Auth{}
	auth.NfcKey = uid
	o := orm.NewOrm()

	// Get user ID
	if err := o.Read(&auth, "NfcKey"); err != nil {
		return "", 0, fmt.Errorf("Failed to read auth table: %v", err)
	}

	// Get user name
	user := User{
		Id: auth.UserId,
	}

	if err := o.Read(&user, "Id"); err != nil {
		return "", 0, fmt.Errorf("Failed to read user table: %v", err)
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
		_, err = o.Insert(&auth)
	} else {
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
		return fmt.Errorf("Failure while checking for auth record: %v", err)
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
	if authRecordMissing {
		beego.Trace("Insert auth")
		_, err = o.Insert(&auth)
	} else {
		beego.Trace("Update auth")
		num, err = o.Update(&auth, "NfcKey")
		if num != 1 {
			return fmt.Errorf("Updated %v rows (expected 1)", num)
		}
	}
	if err != nil {
		return fmt.Errorf("Failed to update: %v", err)
	}
	return nil
}

func DeleteUserAuth(userId int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Auth{UserId: userId})
	return err
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

func createPwResetKey() (string, error) {
	key := make([]byte, PW_RESET_KEY_BYTES)
	_, err := io.ReadFull(rand.Reader, key)
	return fmt.Sprintf("%x", key), err
}

func AuthForgotPassword(email string) (pwResetKey string, err error) {
	email = strings.TrimSpace(email)
	if err = checkEmail(email); err != nil {
		return
	}
	o := orm.NewOrm()
	if pwResetKey, err = createPwResetKey(); err != nil {
		return
	}
	var us []User
	u := User{}
	_, err = o.QueryTable(u.TableName()).
		Filter("email", email).
		All(&us)
	if len(us) != 1 {
		return
	}
	uid := us[0].Id
	a := Auth{UserId: uid}
	if err = o.Read(&a); err != nil {
		return
	}
	a.PwResetKey = pwResetKey
	a.PwResetTime = time.Now()
	n, err := o.Update(&a)
	if err != nil {
		return
	}
	if n == 0 {
		return "", fmt.Errorf("no user with email '%v' found", email)
	}
	return
}

var (
	ErrAuthOutdatedKey = errors.New("Outdated key")
	ErrAuthWrongKey    = errors.New("Wrong key")
	ErrAuthWrongPhone  = errors.New("Wrong phone")
)

func AuthCheckPhone(key string, phone string) (uid int64, err error) {
	if len(key) < PW_RESET_KEY_BYTES/2 {
		return 0, ErrAuthWrongKey
	}
	o := orm.NewOrm()
	a := Auth{}
	var as []Auth
	_, err = o.QueryTable(a.TableName()).
		Filter("pw_reset_key", key).
		All(&as)
	if len(as) == 0 {
		return 0, ErrAuthWrongKey
	} else if len(as) == 1 {
		a = as[0]
		uid = a.UserId
		u, err := GetUser(uid)
		if err != nil {
			return 0, fmt.Errorf("get user %v: %v", uid, err)
		}
		if authPhoneEquals(u.Phone, phone) {
			if a.PwResetTime.After(time.Now()) {
				beego.Info("a.PwResetTime:", a.PwResetTime)
				beego.Info("time.Now():", time.Now())
				return 0, fmt.Errorf("pw reset time is in the future")
			} else if delta := time.Now().Sub(a.PwResetTime); delta > time.Hour {
				beego.Error("key was generated", delta, "ago")
				return 0, ErrAuthOutdatedKey
			} else {
				return uid, nil
			}
		} else {
			return 0, ErrAuthWrongPhone
		}
	} else {
		return 0, fmt.Errorf("two users with same key found!!")
	}
}

func authPhoneEquals(phone1, phone2 string) bool {
	phone1 = authCanonicalizePhone(phone1)
	phone2 = authCanonicalizePhone(phone2)
	return phone1 == phone2
}

func authCanonicalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	removable := []string{
		" ",
		"\n",
		"\t",
		"\r",
		"-",
	}
	for _, ch := range removable {
		phone = strings.Replace(phone, ch, "", -1)
	}
	return phone
}
