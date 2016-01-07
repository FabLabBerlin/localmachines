package modelTest

import (
	"testing"

	"github.com/FabLabBerlin/localmachines/models"
	"github.com/FabLabBerlin/localmachines/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	setup.ConfigDB()
}

func TestMachine(t *testing.T) {
	Convey("Testing Machine model", t, func() {
		Reset(setup.ResetDB)
		Convey("Testing CreateMachine", func() {
			machineName := "My lovely machine"
			Convey("Creating a machine", func() {
				_, err := models.CreateMachine(machineName)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing GetMachine", func() {
			machineName := "My lovely machine"
			Convey("Creating a machine and trying to get it", func() {
				mid, _ := models.CreateMachine(machineName)
				machine, err := models.GetMachine(mid)

				So(machine.Name, ShouldEqual, machineName)
				So(err, ShouldBeNil)
			})
			Convey("Trying to get a non-existing machine should fail", func() {
				_, err := models.GetMachine(0)

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing GetAllMachines", func() {
			machineOneName := "My first machine"
			machineTwoName := "My second lovely machine <3"
			Convey("GetAllMachines when there are no machines in the database", func() {
				machines, err := models.GetAllMachines(false)

				So(len(machines), ShouldEqual, 0)
				So(err, ShouldBeNil)
			})
			Convey("Creating two machines and get them all", func() {
				models.CreateMachine(machineOneName)
				models.CreateMachine(machineTwoName)

				machines, err := models.GetAllMachines(false)

				So(len(machines), ShouldEqual, 2)
				So(err, ShouldBeNil)
			})
		})
		Convey("Testing UpdateMachine", func() {
			machineName := "My lovely machine"
			newMachineName := "This new name is soooooooooooo cool :)"
			Convey("Creating a machine and update it", func() {
				mid, _ := models.CreateMachine(machineName)
				machine, _ := models.GetMachine(mid)
				machine.Name = newMachineName

				err := machine.Update()
				machine, _ = models.GetMachine(mid)
				So(err, ShouldBeNil)
				So(machine.Name, ShouldEqual, newMachineName)
			})
		})
	})
}
