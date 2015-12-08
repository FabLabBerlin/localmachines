package products

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kr15h/fabsmith/models"
)

type Tutor struct {
	Product Product
}

// Get a list of tutors
func GetAllTutors() ([]*Tutor, error) {
	tutorsAsProducts, err := GetAllOfType(TYPE_TUTOR)
	if err != nil {
		msg := "Failed to get tutors as products"
		beego.Error(msg)
		return nil, fmt.Errorf(msg+"%v: ", err)
	}

	tutorList := make([]*Tutor, 0, len(tutorsAsProducts))
	for _, product := range tutorsAsProducts {
		tutor := &Tutor{
			Product: *product,
		}
		tutorList = append(tutorList, tutor)
	}

	return tutorList, nil
}

func CreateTutor(tutor *Tutor) (*Tutor, error) {
	o := orm.NewOrm()

	// In case the type has not been added in upper layers
	tutor.Product.Type = TYPE_TUTOR

	var productId int64
	productId, err := o.Insert(&tutor.Product)
	if err != nil {
		msg := "Failed to insert tutor"
		beego.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	tutor.Product.Id = productId

	return tutor, nil
}

func GetTutor(id int64) (tutor *Tutor, err error) {
	tutor = &Tutor{}
	tutor.Product.Id = id

	o := orm.NewOrm()
	err = o.Read(&tutor.Product)

	return
}

func UpdateTutor(tutor *Tutor) error {
	if tutor.Product.UserId != 0 {
		o := orm.NewOrm()
		// Get user name by user ID
		user := models.User{}
		user.Id = tutor.Product.UserId
		err := o.Read(&user)
		if err != nil {
			beego.Error("Failed to read user:", err)
			return fmt.Errorf("Failed to read user: %v", err)
		}
		tutor.Product.Name = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}
	return Update(&tutor.Product)
}
