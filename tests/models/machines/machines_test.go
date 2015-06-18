package machineTests

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

func TestDeleteMachine(t *testing.T) {
	Convey("Testing DeleteMachine", t, func() {
		machineName := "My lovely machine"
		Convey("Creating a machine and delete it", func() {
			mid, _ := models.CreateMachine(machineName)

			err := models.DeleteMachine(mid)
			So(err, ShouldBeNil)
		})
		Convey("Try to delete non-existing user", func() {
			err := models.DeleteMachine(0)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestCreateMachine(t *testing.T) {
	Convey("Testing CreateMachine", t, func() {
		machineName := "My lovely machine"
		Convey("Creating a machine", func() {
			mid, err := models.CreateMachine(machineName)
			defer models.DeleteMachine(mid)

			So(err, ShouldBeNil)
		})
	})
}

func TestGetMachine(t *testing.T) {
	Convey("Testing GetMachine", t, func() {
		machineName := "My lovely machine"
		Convey("Creating a machine and trying to get it", func() {
			mid, _ := models.CreateMachine(machineName)
			defer models.DeleteMachine(mid)
			machine, err := models.GetMachine(mid)

			So(machine.Name, ShouldEqual, machineName)
			So(err, ShouldBeNil)
		})
		Convey("Trying to get a non-existing machine", func() {
			machine, err := models.GetMachine(0)

			So(machine, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetAllMachines(t *testing.T) {
	Convey("Testing GetAllMachines", t, func() {
		machineOneName := "My first machine"
		machineTwoName := "My second lovely machine <3"
		Convey("GetAllMachines when there are no machines in the database", func() {
			machines, err := models.GetAllMachines()

			So(len(machines), ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
		Convey("Creating two machines and get them all", func() {
			mid1, _ := models.CreateMachine(machineOneName)
			defer models.DeleteMachine(mid1)
			mid2, _ := models.CreateMachine(machineTwoName)
			defer models.DeleteMachine(mid2)

			machines, err := models.GetAllMachines()

			So(len(machines), ShouldEqual, 2)
			So(err, ShouldBeNil)
		})
	})
}

func TestUpdateMachine(t *testing.T) {
	Convey("Testing UpdateMachine", t, func() {
		machineName := "My lovely machine"
		newMachineName := "This new name is soooooooooooo cool :)"
		Convey("Creating a machine and update it", func() {
			mid, _ := models.CreateMachine(machineName)
			defer models.DeleteMachine(mid)
			machine, _ := models.GetMachine(mid)
			machine.Name = newMachineName

			err := models.UpdateMachine(machine)
			machine, _ = models.GetMachine(mid)
			So(err, ShouldBeNil)
			So(machine.Name, ShouldEqual, newMachineName)
		})
	})
}
