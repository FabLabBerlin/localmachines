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
				mid, _ := machine.Create(1, machineName)
				machine, err := machine.Get(mid)

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
				mid, _ := machine.Create(1, machineName)
				m, _ := machine.Get(mid)
				m.LocationId = 1
				m.Name = newMachineName

				err := m.Update()
				m, _ = machine.Get(mid)
				So(err, ShouldBeNil)
				So(m.Name, ShouldEqual, newMachineName)
			})
		})
		Convey("Netswitch url on/off 200 with should give no error", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			ns := machine.Machine{
				NetswitchUrlOn:  ts.URL + "?method=on",
				NetswitchUrlOff: ts.URL + "?method=off",
				NetswitchXmpp:   false,
			}
			So(ns.On(), ShouldBeNil)
			So(ns.Off(), ShouldBeNil)
		})
		Convey("Netswitch url on/off 500 with should give an error", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer ts.Close()

			ns := machine.Machine{
				NetswitchUrlOn:  ts.URL + "?method=on",
				NetswitchUrlOff: ts.URL + "?method=off",
				NetswitchXmpp:   false,
			}
			So(ns.On(), ShouldNotBeNil)
			So(ns.Off(), ShouldNotBeNil)
		})
		Convey("Creating same netswitch mapping for two machines", func() {
			mid1, err := machine.Create(1, "foo")
			if err != nil {
				panic(err.Error())
			}
			mid2, err := machine.Create(1, "bar")
			if err != nil {
				panic(err.Error())
			}
			m1, err := machine.Get(mid1)
			if err != nil {
				panic(err.Error())
			}
			m2, err := machine.Get(mid2)
			if err != nil {
				panic(err.Error())
			}
			m1.NetswitchHost = "example.com"
			m1.NetswitchXmpp = true
			err = m1.Update()
			So(err, ShouldBeNil)
			m2.NetswitchHost = "example.com"
			m2.NetswitchXmpp = true
			err = m2.Update()
			So(err, ShouldNotBeNil)
		})
	})
}
