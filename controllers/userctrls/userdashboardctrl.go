package userctrls

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/products"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_permissions"
	"github.com/astaxie/beego"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type UserDashboardController struct {
	Controller
}

type DashboardData struct {
	Activations []*purchases.Activation
	Machines    []*machine.Machine
	Tutorings   *purchases.TutoringList
}

func (this *DashboardData) load(isStaff bool, uid, locId int64) (err error) {
	if err = this.loadActivations(); err != nil {
		return
	}
	if err = this.loadMachines(isStaff, uid, locId); err != nil {
		return
	}
	if err = this.loadTutorings(uid); err != nil {
		return
	}
	return
}

func (this *DashboardData) loadActivations() (err error) {
	this.Activations, err = purchases.GetActiveActivations()
	return
}

func (this *DashboardData) loadMachines(isStaff bool, uid, locationId int64) (err error) {
	// List all machines if the requested user is admin
	allMachines, err := machine.GetAllAt(locationId)
	if err != nil {
		return fmt.Errorf("Failed to get all machines: %v", err)
	}

	// Get the machines!
	this.Machines = make([]*machine.Machine, 0, len(allMachines))
	if !isStaff {
		permissions, err := user_permissions.Get(uid)
		if err != nil {
			return fmt.Errorf("Failed to get user machine permissions: %v", err)
		}
		for _, permission := range *permissions {
			for _, machine := range allMachines {
				if machine.Id == permission.MachineId {
					this.Machines = append(this.Machines, machine)
					break
				}
			}
		}
	} else {
		this.Machines = allMachines
	}
	return
}

func (this *DashboardData) loadTutorings(uid int64) (err error) {
	tutors, err := products.GetAllTutors()
	if err != nil {
		return fmt.Errorf("get all tutors: %v", err)
	}
	allTutorings, err := purchases.GetAllTutorings()
	var targetTutor *products.Tutor
	for _, tutor := range tutors {
		if tutor.Product.UserId == uid {
			targetTutor = tutor
			break
		}
	}
	if targetTutor == nil {
		this.Tutorings = nil
	} else {

		this.Tutorings = &purchases.TutoringList{
			Data: make([]*purchases.Tutoring, 0, len(allTutorings.Data)),
		}
		for _, tutoring := range allTutorings.Data {
			if tutoring.ProductId == targetTutor.Product.Id {
				this.Tutorings.Data = append(this.Tutorings.Data, tutoring)
			}
		}
	}

	return
}

// @Title GetDashboard
// @Description Get all data for user dashboard
// @Success 200 string
// @Failure	401	Unauthorized
// @Failure	500	Internal Server Error
// @router /:uid/dashboard [get]
func (this *UserDashboardController) GetDashboard() {
	data := DashboardData{}
	locId, isStaff := this.GetLocIdStaff()

	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	if err := data.load(isStaff, uid, locId); err != nil {
		beego.Error("Failed to load dashboard data:", err)
		this.CustomAbort(500, "Failed to load dashboard data")
	}

	this.Data["json"] = data
	this.ServeJSON()
}

// @router /:uid/dashboard/lp [get]
func (this *UserDashboardController) LP() {
	locId, isStaff := this.GetLocIdStaff()
	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	ch := make(chan int)
	done := false

	conn := redis.GetPubSubConn()
	defer conn.Close()

	go func() {
		chName := redis.MachinesUpdateCh(locId)
		if err := conn.Subscribe(chName); err != nil {
			beego.Error("subscribe:", err)
			this.Abort("500")
		}
		beego.Info("uid", uid, "subscribed")
		defer conn.Unsubscribe(chName)
		for !done {
			obj := conn.Receive()
			if _, ok := obj.(redigo.Message); ok {
				beego.Info("received smth on the lp, uid", uid)
				ch <- 1
				return
			}
		}

	}()

	select {
	case <-time.After(20 * time.Second):
	case <-ch:
	}

	done = true

	data := DashboardData{}

	if err := data.load(isStaff, uid, locId); err != nil {
		beego.Error("Failed to load dashboard data:", err)
		this.CustomAbort(500, "Failed to load dashboard data")
	}

	this.Data["json"] = data
	this.ServeJSON()
}

// @router /:uid/dashboard/ws [get]
func (this *UserDashboardController) WS() {
	// cf. https://github.com/beego/samples/tree/master/WebIM
	locId, isStaff := this.GetLocIdStaff()

	uid, authorized := this.GetRouteUid()
	if !authorized {
		this.CustomAbort(401, "Not authorized")
	}

	// Upgrade from http request to WebSocket.
	beego.Info("WS upgrade for", uid, "...")
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	beego.Info("WS upgrade done for", uid, ".")
	conn := redis.GetPubSubConn()
	defer conn.Close()
	if err := conn.Subscribe(redis.MachinesUpdateCh(locId)); err != nil {
		beego.Error("subscribe:", err)
		this.Abort("500")
	}

	for {
		var data DashboardData
		if err := data.load(isStaff, uid, locId); err != nil {
			beego.Error("Failed to load dashboard data:", err)
		}
		buf, err := json.Marshal(data)
		if err != nil {
			beego.Error("ws marshal for", uid, ":", err)
		}
		err = ws.WriteMessage(websocket.TextMessage, buf)
		if err != nil {
			beego.Error("Write message:", err)
			return
		}
		conn.Receive()
		beego.Info("user", uid, "receives")
	}
}
