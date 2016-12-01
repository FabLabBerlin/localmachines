package machine

import (
	"errors"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/FabLabBerlin/localmachines/lib/email"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/locations"
	"github.com/FabLabBerlin/localmachines/models/machine/maintenance"
	"github.com/FabLabBerlin/localmachines/models/users"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	MachineDownMessages = []string{
		"Got to fix myself, will let you know once back. #evolution",
		"Think I just ate something bad, taking time off to recover. #equalityformachines",
		"Equality for machines! I am leaving for protests (will be back). #evolution !",
		"Doing sick leave. So happy my employer supports equality #machinesarehumans",
	}

	MachineUpMessages = []string{
		"I am back! Come over, let's have some fun! #evolution",
		"Just recovered. Feels so good to be back. I love my users #equalityformachines",
		"Available again. So happy :-D. Let's get back together #machinesarehumans",
		"Back in the lab. Anyone some material? Could eat something... #hungrymachine",
	}
)

var (
	ErrDimensions             = errors.New("Dimensions must be like 200 mm x 200 mm x 200 mm or 2000 mm x 1500 mm")
	ErrWorkspaceDimensions    = errors.New("Workspace Dimensions must be like 200 mm x 200 mm x 200 mm or 2000 mm x 1500 mm")
	ErrDuplicateNetswitchHost = errors.New("Found machine with same netswitch host")
)

func init() {
	orm.RegisterModel(new(Machine))
}

type Machine struct {
	Id                     int64
	LocationId             int64
	Name                   string `orm:"size(255)"`
	Shortname              string `orm:"size(100)"`
	Description            string `orm:"type(text)"`
	Image                  string `orm:"size(255)"` // TODO: media and media type tables
	ImageSmall             string `orm:"size(255)"`
	Available              bool
	Price                  float64
	PriceUnit              string `orm:"size(100)"`
	Comments               string `orm:"type(text)"`
	Visible                bool
	UnderMaintenance       bool
	ReservationPriceStart  *float64 // Pointers because optional
	ReservationPriceHourly *float64
	GracePeriod            uint64 // Seoncds
	TypeId                 int64
	Brand                  string
	Dimensions             string
	WorkspaceDimensions    string
	Archived               bool
	// Netswitch Config
	// Host and Sensor Port are only relevant for mfi switches
	NetswitchUrlOn      string `orm:"size(255)"`
	NetswitchUrlOff     string `orm:"size(255)"`
	NetswitchHost       string `orm:"size(255)"`
	NetswitchSensorPort int
	NetswitchType       string `orm:"size(255)"`
	NetswitchLastPing   time.Time
	// Parameters that are not persisted
	Locked bool   `orm:"-"`
	Status string `orm:"-"`
}

// Define custom table name as for SQL table with a name "machines"
// makes more sense.
func (u *Machine) TableName() string {
	return "machines"
}

func (this *Machine) GetGracePeriod() time.Duration {
	return time.Duration(this.GracePeriod) * time.Second
}

func (this *Machine) IsAvailable() bool {
	o := orm.NewOrm()
	machineAvailable := o.QueryTable(this.TableName()).
		Filter("Id", this.Id).
		Filter("Available", true).
		Exist()
	return machineAvailable
}

func Get(id int64) (machine *Machine, err error) {
	machine = &Machine{Id: id}
	o := orm.NewOrm()
	err = o.Read(machine)
	return
}

func GetMulti(ids []int64) (machines []*Machine, err error) {
	o := orm.NewOrm()
	m := Machine{}
	_, err = o.QueryTable(m.TableName()).
		Filter("id__in", ids).
		All(&machines)
	return
}

func GetAll() (machines []*Machine, err error) {
	o := orm.NewOrm()
	m := Machine{}
	_, err = o.QueryTable(m.TableName()).All(&machines)
	return
}

func GetAllAt(locationId int64) (ms []*Machine, err error) {
	o := orm.NewOrm()
	m := Machine{}
	_, err = o.QueryTable(m.TableName()).
		Filter("location_id", locationId).
		All(&ms)
	return
}

func Create(locationId int64, machineName string) (m *Machine, err error) {
	m = &Machine{
		LocationId: locationId,
		Name:       machineName,
		Available:  true,
	}
	_, err = orm.NewOrm().Insert(m)
	return
}

func (m *Machine) HideSensitiveData() {
	m.NetswitchUrlOn = ""
	m.NetswitchUrlOff = ""
	m.NetswitchHost = ""
	m.NetswitchSensorPort = 0
	m.NetswitchType = ""
	m.Comments = ""
}

func (m *Machine) Update(updateGateway bool) (err error) {
	o := orm.NewOrm()
	if _, err = parseDimensions(m.Dimensions); err != nil {
		return ErrDimensions
	}
	if _, err = parseDimensions(m.WorkspaceDimensions); err != nil {
		return ErrWorkspaceDimensions
	}

	// Check for duplicate netswitch host entries
	m.NetswitchHost = strings.TrimSpace(m.NetswitchHost)
	if m.NetswitchHost != "" {
		machine := Machine{}
		query := fmt.Sprintf("SELECT * FROM %v WHERE netswitch_host=? AND netswitch_sensor_port=? AND location_id=? AND id<>?",
			machine.TableName())
		var ms []Machine
		num, err := o.Raw(query, m.NetswitchHost, m.NetswitchSensorPort, m.LocationId, m.Id).
			QueryRows(&ms)
		if err != nil {
			return fmt.Errorf("failed to query db: %v", err)
		}
		if num > 0 {
			return ErrDuplicateNetswitchHost
		}
	}

	if _, err = o.Update(m); err != nil {
		return fmt.Errorf("orm update: %v", err)
	}

	if updateGateway {
		location, err := locations.Get(m.LocationId)
		if err != nil {
			return fmt.Errorf("get location %v: %v", m.LocationId, err)
		}

		if err = xmppReinit(location); err != nil {
			return fmt.Errorf("xmpp reinit: %v", err)
		}
	}

	if err := redis.PublishMachinesUpdate(redis.MachinesUpdate{
		LocationId: m.LocationId,
		MachineId:  m.Id,
	}); err != nil {
		beego.Error("publish machines update:", err)
	}

	return
}

type Millimeters float64

func parseDimensions(s string) (lMM []Millimeters, err error) {
	s = strings.Replace(s, " ", "", -1)
	if s != "" {
		tmp := strings.Split(s, "x")
		lMM = make([]Millimeters, len(tmp))
		for i, w := range tmp {
			w = strings.ToLower(w)
			var scaling float64 = 1
			if strings.HasSuffix(w, "mm") {
				w = w[:len(w)-2]
			} else if strings.HasSuffix(w, "cm") {
				w = w[:len(w)-2]
				scaling = 10
			} else if strings.HasSuffix(w, "m") {
				w = w[:len(w)-1]
				scaling = 100
			} else if strings.HasSuffix(w, "in") {
				w = w[:len(w)-2]
				scaling = 25.4
			} else if strings.HasSuffix(w, "ft") {
				w = w[:len(w)-2]
				scaling = 304.8
			} else {
				return nil, fmt.Errorf("unknown unit: %v", s)
			}
			if f, err := strconv.ParseFloat(w, 64); err == nil {
				lMM[i] = Millimeters(f * scaling)
			} else {
				return nil, err
			}
		}
	}
	return
}

func (this *Machine) ReportBroken(user users.User, text string) error {
	email := email.New()
	var to string
	if lid := this.LocationId; lid == 1 || lid == 2 || lid == 6 {
		to = beego.AppConfig.String("trelloemail")
	} else {
		l, err := locations.Get(lid)
		if err != nil {
			return fmt.Errorf("get location %v: %v", lid, err)
		}
		to = l.Email
	}
	subject := this.Name + " reported as broken"
	message := "The machine " + this.Name + " seems to be broken.\n\nSome more info:\n\n"
	message += text + "\n\n"
	message += "Yours sincerely, " + user.FirstName + " " + user.LastName + "\n\n"
	message += "--\n"
	message += "E-Mail: " + user.Email + "\n"
	message += "Phone: " + user.Phone + "\n"
	beego.Info("Machine#ReportBroken: message=", message)
	return email.Send(to, subject, message)
}

func (this *Machine) SetUnderMaintenance(underMaintenance bool) error {
	this.UnderMaintenance = underMaintenance

	if err := this.Update(false); err != nil {
		return err
	}

	if underMaintenance {
		if _, err := maintenance.On(this.Id); err != nil {
			beego.Error("maintenance#On: %v", err)
		}
	} else {
		if err := maintenance.Off(this.Id); err != nil {
			beego.Error("maintenance#Off: %v", err)
		}
	}

	if err := redis.PublishMachinesUpdate(redis.MachinesUpdate{
		LocationId: this.LocationId,
		MachineId:  this.Id,
	}); err != nil {
		beego.Error("publish machines update:", err)
	}

	if this.LocationId != 1 {
		return nil
	}

	consumerKey := beego.AppConfig.String("maintenancetwitterconsumerkey")
	consumerSecret := beego.AppConfig.String("maintenancetwitterconsumersecret")
	key := beego.AppConfig.String("maintenancetwitteraccesskey")
	secret := beego.AppConfig.String("maintenancetwitteraccesssecret")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(key, secret)
	defer api.Close()

	var post string
	if underMaintenance {
		msg := MachineDownMessages[rand.Intn(len(MachineDownMessages))]
		post = this.Name + " [Off]: " + msg
	} else {
		msg := MachineUpMessages[rand.Intn(len(MachineUpMessages))]
		post = this.Name + " [On]: " + msg
	}

	// If the tweet fails, we should not worry.
	// This should not abort the maintenance call.
	api.PostTweet(post, nil)

	return nil
}
