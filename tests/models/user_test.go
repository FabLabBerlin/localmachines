package tests

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"github.com/kr15h/fabsmith/models"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	// find app.conf path
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator)+".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

	// Setting log level
	beego.SetLevel(beego.LevelError)

	// Force Runmode to "test"
	beego.RunMode = "test"

	// Get MySQL config from environment variables
	mysqlUser := beego.AppConfig.String("mysqluser")
	if mysqlUser == "" {
		panic("Please set mysqluser in app.conf")
	}

	mysqlPass := beego.AppConfig.String("mysqlpass")
	if mysqlPass == "" {
		panic("Please set mysqlpass in app.conf")
	}

	mysqlHost := beego.AppConfig.String("mysqlhost")
	if mysqlHost == "" {
		mysqlHost = "localhost"
	}

	mysqlPort := beego.AppConfig.String("mysqlport")
	if mysqlPort == "" {
		mysqlPort = "3306"
	}

	mysqlDb := beego.AppConfig.String("mysqldb")
	beego.Info("Lel: " + mysqlDb)
	if mysqlDb == "" {
		panic("Please set mysqldb in app.conf")
	}

	// Build MySQL connection string out of the config variables
	mysqlConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDb)

	// Register MySQL driver and default database for beego ORM
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", mysqlConnString)
}

func TestUserCreate(t *testing.T) {
	Convey("Testing CreateUser", t, func() {
		u := models.User{
			FirstName: "test",
			LastName:  "test",
		}
		Convey("Creating one user into database", func() {
			// Creating user
			uc, err := models.CreateUser(&u)
			So(err, ShouldBeNil)
			So(uc, ShouldBeGreaterThan, 0)

			// Deleting user for nexts tests
			models.DeleteUser(uc)
		})
		Convey("Creating 2 users that are identical into database, should get an error", func() {
			// Creating first user
			uc, err := models.CreateUser(&u)
			So(err, ShouldBeNil)
			So(uc, ShouldBeGreaterThan, 0)

			// Creating second user
			uc2, err2 := models.CreateUser(&u)
			So(err2, ShouldNotBeNil)
			So(uc2, ShouldEqual, 0)

			// Deleting user for nexts tests
			models.DeleteUser(uc)
		})
	})
}

func TestAuthenticateUser(t *testing.T) {
	Convey("Testing AuthenticateUser", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating a user and setting him a password", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			So(err, ShouldBeNil)
			So(uid, ShouldBeGreaterThan, 0)
			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)
			So(err, ShouldBeNil)
		})
		Convey("Creating a user with a password and try to authenticate him", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			So(err, ShouldBeNil)
			So(uid, ShouldBeGreaterThan, 0)
			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)
			So(err, ShouldBeNil)

			authUID, err := models.AuthenticateUser("test", "test")
			So(authUID, ShouldEqual, uid)
			So(err, ShouldBeNil)
		})
		Convey("Creating a user with a password and try to authenticate with wrong username", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			So(err, ShouldBeNil)
			So(uid, ShouldBeGreaterThan, 0)
			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)
			So(err, ShouldBeNil)

			authUID, err := models.AuthenticateUser("wrong", "test")
			So(authUID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
		Convey("Creating a user with a password and try to authenticate with wrong password", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			So(err, ShouldBeNil)
			So(uid, ShouldBeGreaterThan, 0)
			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)
			So(err, ShouldBeNil)

			authUID, err := models.AuthenticateUser("test", "wrong")
			So(authUID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAuthenticateUserUID(t *testing.T) {
	Convey("Testing AuthenticateUser", t, func() {
		u := models.User{
			Username: "test",
		}
		Convey("Creating a user and setting him a NFC UID", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			So(err, ShouldBeNil)
			So(uid, ShouldBeGreaterThan, 0)
			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)
			So(err, ShouldBeNil)
			err = models.AuthUpdateNfcUid(uid, "123456")
			So(err, ShouldBeNil)
		})
		Convey("Creating a user with a NFC UID and try to authenticate him", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			So(err, ShouldBeNil)
			So(uid, ShouldBeGreaterThan, 0)
			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)
			So(err, ShouldBeNil)
			err = models.AuthUpdateNfcUid(uid, "123456")
			So(err, ShouldBeNil)
			_, authUID, err := models.AuthenticateUserUid("123456")
			So(authUID, ShouldEqual, uid)
			So(err, ShouldBeNil)
		})
		Convey("Creating a user with a NFC UID and try to authenticate him with wrong UID", func() {
			uid, err := models.CreateUser(&u)
			defer models.DeleteUser(uid)

			So(err, ShouldBeNil)
			So(uid, ShouldBeGreaterThan, 0)
			err = models.AuthSetPassword(uid, "test")
			defer models.DeleteUserAuth(uid)
			So(err, ShouldBeNil)
			err = models.AuthUpdateNfcUid(uid, "123456")
			So(err, ShouldBeNil)
			_, authUID, err := models.AuthenticateUserUid("654321")
			So(authUID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
	})
}
