package membershipTests

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
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator)+".."+string(filepath.Separator)+".."+string(filepath.Separator))))
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

func TestGetAllMemberships(t *testing.T) {
	Convey("Testing GetAllMemberships", t, func() {
		membershipName := "Lel"
		Convey("Getting all memberships when there is nothing in the database", func() {
			memberships, err := models.GetAllMemberships()

			So(err, ShouldBeNil)
			So(len(memberships), ShouldEqual, 0)
		})
		Convey("Creating 2 memberships and get them all", func() {
			mid1, _ := models.CreateMembership(membershipName)
			defer models.DeleteMembership(mid1)
			mid2, _ := models.CreateMembership(membershipName)
			defer models.DeleteMembership(mid2)

			memberships, err := models.GetAllMemberships()

			So(err, ShouldBeNil)
			So(len(memberships), ShouldEqual, 2)
		})
	})
}

func TestCreateMembership(t *testing.T) {
	Convey("Testing CreateMembership", t, func() {
		membershipName := "My awesome membership"
		Convey("Creating one membership", func() {
			mid, err := models.CreateMembership(membershipName)
			defer models.DeleteMembership(mid)

			So(err, ShouldBeNil)
		})
	})
}

func TestGetMembership(t *testing.T) {
	Convey("Testing GetMembership", t, func() {
		membershipName := "Lel"
		Convey("Getting non-existing membership", func() {
			_, err := models.GetMembership(0)

			So(err, ShouldNotBeNil)
		})
		Convey("Creating a membership and getting it", func() {
			mid, _ := models.CreateMembership(membershipName)
			defer models.DeleteMembership(mid)
			membership, err := models.GetMembership(mid)

			So(membership.Id, ShouldEqual, mid)
			So(err, ShouldBeNil)
		})
	})
}

func TestUpdateMembership(t *testing.T) {
	Convey("Testing UpdateMembership", t, func() {
		membershipName := "Lel"
		newMembershipName := "DatAwesomeNewName"
		Convey("Try updating with nil object", func() {
			panicFunc := func() {
				models.UpdateMembership(nil)
			}

			So(panicFunc, ShouldPanic)
		})
		Convey("Try to update non existing membership", func() {
			m := &models.Membership{
				Title: membershipName,
			}

			err := models.UpdateMembership(m)
			So(err, ShouldNotBeNil)
		})
		Convey("Create membership and update it", func() {
			mid, _ := models.CreateMembership(membershipName)
			defer models.DeleteMembership(mid)

			m, _ := models.GetMembership(mid)
			m.Title = newMembershipName

			err := models.UpdateMembership(m)
			nm, _ := models.GetMembership(mid)
			So(err, ShouldBeNil)
			So(nm.Title, ShouldEqual, newMembershipName)
		})
	})
}

func TestDeleteMembership(t *testing.T) {
	Convey("Testing DeleteMembership", t, func() {
		membershipName := "My awesome membership program"
		Convey("Try to delete non-existing membership", func() {
			err := models.DeleteMembership(0)

			So(err, ShouldNotBeNil)
		})
		Convey("Creating a membership and delete it", func() {
			mid, _ := models.CreateMembership(membershipName)
			err := models.DeleteMembership(mid)

			So(err, ShouldBeNil)
		})
	})
}

func TestCreateUserMembership(t *testing.T) {
	Convey("Testing CreateUserMembership", t, func() {
		Convey("Try creating a user with nil parameter", func() {
			_, err := models.CreateUserMembership(nil)

			So(err, ShouldNotBeNil)
		})
	})
}
