package modelTest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FabLabBerlin/localmachines/models/machine"
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
				_, err := machine.Create(1, machineName)

				So(err, ShouldBeNil)
			})
		})
		Convey("Testing GetMachine", func() {
			machineName := "My lovely machine"
			Convey("Creating a machine and trying to get it", func() {
				m, _ := machine.Create(1, machineName)
				machine, err := machine.Get(m.Id)

				So(machine.Name, ShouldEqual, machineName)
				So(err, ShouldBeNil)
			})
			Convey("Trying to get a non-existing machine should fail", func() {
				_, err := machine.Get(0)

				So(err, ShouldNotBeNil)
			})
		})
		Convey("Testing GetAllMachines", func() {
			machineOneName := "My first machine"
			machineTwoName := "My second lovely machine <3"
			Convey("GetAllMachines when there are no machines in the database", func() {
				machines, err := machine.GetAll()

				So(len(machines), ShouldEqual, 0)
				So(err, ShouldBeNil)
			})
			Convey("Creating two machines and get them all", func() {
				machine.Create(1, machineOneName)
				machine.Create(1, machineTwoName)

				machines, err := machine.GetAll()

				So(len(machines), ShouldEqual, 2)
				So(err, ShouldBeNil)
			})
		})
		Convey("Testing UpdateMachine", func() {
			machineName := "My lovely machine"
			newMachineName := "This new name is soooooooooooo cool :)"
			Convey("Creating a machine and update it", func() {
				m, _ := machine.Create(1, machineName)
				m.LocationId = 1
				m.Name = newMachineName

				err := m.Update(false)
				m, _ = machine.Get(m.Id)
				So(err, ShouldBeNil)
				So(m.Name, ShouldEqual, newMachineName)
			})
		})
		Convey("Netswitch url on/off 500 with should give an error", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer ts.Close()

			ns := machine.Machine{
				NetswitchUrlOn:  ts.URL + "?method=on",
				NetswitchUrlOff: ts.URL + "?method=off",
			}
			So(ns.On(0), ShouldNotBeNil)
			So(ns.Off(0), ShouldNotBeNil)
		})
		Convey("Creating same netswitch mapping for two machines", func() {
			m1, err := machine.Create(1, "foo")
			if err != nil {
				panic(err.Error())
			}
			m2, err := machine.Create(1, "bar")
			if err != nil {
				panic(err.Error())
			}
			m1.NetswitchHost = "example.com"
			err = m1.Update(false)
			So(err, ShouldBeNil)
			m2.NetswitchHost = "example.com"
			err = m2.Update(false)
			So(err, ShouldNotBeNil)
		})
		Convey("Creating same netswitch mapping for two machines at different locations", func() {
			m1, err := machine.Create(1, "foo")
			if err != nil {
				panic(err.Error())
			}
			m2, err := machine.Create(2, "bar")
			if err != nil {
				panic(err.Error())
			}
			m1.NetswitchHost = "example.com"
			err = m1.Update(false)
			So(err, ShouldBeNil)
			m2.NetswitchHost = "example.com"
			err = m2.Update(false)
			So(err, ShouldBeNil)
		})
	})
}
