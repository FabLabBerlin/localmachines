package userctrls

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/products"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/user_permissions"
	"github.com/astaxie/beego"
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
