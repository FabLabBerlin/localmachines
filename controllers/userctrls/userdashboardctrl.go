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

const WS_PING_INTERVAL_SECONDS = 29

type UserDashboardController struct {
	Controller
}

type PushData struct {
	Activations []*purchases.Activation
	Machines    []*machine.Machine
	Tutorings   *purchases.TutoringList
	UserMessage struct {
		Error   string
		Info    string
		Warning string
	}
}

func (this *PushData) load(isStaff bool, uid, locId int64) (err error) {
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

func (this *PushData) loadActivations() (err error) {
	this.Activations, err = purchases.GetActiveActivations()
	return
}

func (this *PushData) loadMachines(isStaff bool, uid, locationId int64) (err error) {
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
		for _, machine := range allMachines {
			machine.Locked = true
			for _, permission := range *permissions {
				if machine.TypeId == permission.CategoryId {
					machine.Locked = false
					break
				}
			}
			this.Machines = append(this.Machines, machine)
		}
	} else {
		this.Machines = allMachines
	}

	for _, m := range this.Machines {
		m.HideSensitiveData()
	}

	return
}

func (this *PushData) loadTutorings(uid int64) (err error) {
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
	data := PushData{}
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

	conn := redis.GetPoolConn()
	defer conn.Close()
	psc := &redigo.PubSubConn{
		Conn: conn,
	}
	defer func() {
		psc.Conn.Send("QUIT")
		psc.Conn.Flush()
	}()

	data := PushData{}

	go func() {
		chName := redis.MachinesUpdateCh(locId)
		if err := psc.Subscribe(chName); err != nil {
			beego.Error("subscribe:", err)
			this.Abort("500")
		}
		beego.Info("uid", uid, "subscribed")
		defer psc.Unsubscribe(chName)
		for !done {
			obj := psc.Receive()
			beego.Info("LP: psc.Receive :)")
			if v, ok := obj.(redigo.Message); ok {
				var update redis.MachinesUpdate
				if err := json.Unmarshal(v.Data, &update); err == nil {
					beego.Info("user", uid, "get update", update)
				} else {
					beego.Info("json unmarshal:", err)
				}
				beego.Info("update.Command=", update.Command)
				beego.Info("unmarshal the update=", update)
				if update.UserId == uid {
					data.UserMessage.Error = update.Error
					data.UserMessage.Info = update.Info
					data.UserMessage.Warning = update.Warning
				}
				ch <- 1
				return
			} else {
				beego.Info("it's not a message but some other shit", obj)
			}
		}

	}()

	select {
	case <-time.After(20 * time.Second):
	case <-ch:
	}

	done = true

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

	done := false

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
	defer ws.Close()

	beego.Info("WS upgrade done for", uid, ".")
	conn := redis.GetPoolConn()
	defer conn.Close()
	psc := &redigo.PubSubConn{
		Conn: conn,
	}
	defer func() {
		psc.Conn.Send("QUIT")
		psc.Conn.Flush()
	}()

	chName := redis.MachinesUpdateCh(locId)
	if err := psc.Subscribe(chName); err != nil {
		beego.Error("subscribe:", err)
		this.Abort("500")
	}
	defer psc.Unsubscribe(chName)
	defer beego.Info("WS() END")

	go func() {
		for !done {
			beego.Info("Pinnggg ws...")
			<-time.After(WS_PING_INTERVAL_SECONDS * time.Second)
			err = ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				done = true
				beego.Error("Write ping message:", err)
				return
			}
		}
	}()

	for firstIteration := true; !done; firstIteration = false {
		var data PushData

		if err := data.load(isStaff, uid, locId); err != nil {
			beego.Error("Failed to load dashboard data:", err)
		}

		if !firstIteration {
			msg := psc.Receive()
			beego.Info("psc.Receive()")

			switch v := msg.(type) {
			case redigo.Message:
				var update redis.MachinesUpdate
				if err := json.Unmarshal(v.Data, &update); err == nil {
					beego.Info("user", uid, "get update", update)
				} else {
					beego.Info("json unmarshal:", msg)
				}
				if update.UserId == uid {
					data.UserMessage.Error = update.Error
					data.UserMessage.Info = update.Info
					data.UserMessage.Warning = update.Warning
				}
			default:
				beego.Info("user", uid, "receives", msg)
			}
		}

		buf, err := json.Marshal(data)
		if err != nil {
			beego.Error("ws marshal for", uid, ":", err)
		}
		err = ws.WriteMessage(websocket.TextMessage, buf)
		if err != nil {
			done = true
			beego.Error("Write text message:", err)
			return
		}
	}

	this.Finish()
}
