package userctrls

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/kr15h/fabsmith/controllers"
	"github.com/kr15h/fabsmith/models"
	"github.com/kr15h/fabsmith/models/products"
	"github.com/kr15h/fabsmith/models/purchases"
)

type UserDashboardController struct {
	controllers.Controller
}

type DashboardData struct {
	Activations []*purchases.Activation
	Machines    []*models.Machine
	Tutorings   *purchases.TutoringPurchaseList
}

func (this *DashboardData) load(isAdmin bool, uid int64) (err error) {
	if err = this.loadActivations(); err != nil {
		return
	}
	if err = this.loadMachines(isAdmin, uid); err != nil {
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

func (this *DashboardData) loadMachines(isAdmin bool, uid int64) (err error) {
	// List all machines if the requested user is admin
	allMachines, err := models.GetAllMachines(false)
	if err != nil {
		return fmt.Errorf("Failed to get all machines: %v", err)
	}

	// Get the machines!
	this.Machines = make([]*models.Machine, 0, len(allMachines))
	if !isAdmin {
		permissions, err := models.GetUserPermissions(uid)
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
	allTutorings, err := purchases.GetAllTutoringPurchases()
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

		this.Tutorings = &purchases.TutoringPurchaseList{
			Data: make([]*purchases.TutoringPurchase, 0, len(allTutorings.Data)),
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

	// Check if logged in
	suid := this.GetSession(controllers.SESSION_FIELD_NAME_USER_ID)
	if suid == nil {
		beego.Info("Not logged in")
		this.CustomAbort(401, "Unauthorized")
	}

	// Get requested user ID
	var err error
	var ruid int64
	ruid, err = this.GetInt64(":uid")
	if err != nil {
		beego.Error("Failed to get :uid", err)
		this.CustomAbort(500, "Internal Server Error")
	}

	suidInt64, ok := suid.(int64)
	if !ok {
		beego.Error("Could not get session user ID as int64")
		this.CustomAbort(500, "Internal Server Error")
	}

	if suidInt64 != ruid {
		if !this.IsAdmin() {
			beego.Error("Not authorized")
			this.CustomAbort(401, "Unauthorized")
		}
	}

	if err := data.load(this.IsAdmin(ruid), ruid); err != nil {
		beego.Error("Failed to load dashboard data:", err)
		this.CustomAbort(500, "Failed to load dashboard data")
	}

	this.Data["json"] = data
	this.ServeJson()
}
