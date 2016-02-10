package products

import (
	"github.com/astaxie/beego/orm"
)

type CoWorking struct {
	Product Product
}

func CreateCoWorking(locationId int64, name string) (cp CoWorking, err error) {
	cp = CoWorking{
		Product: Product{
			LocationId: locationId,
			Name:       name,
			Type:       TYPE_CO_WORKING,
		},
	}
	id, err := Create(&cp.Product)
	cp.Product.Id = id
	return
}

func (cp *CoWorking) Update() (err error) {
	return cp.Product.Update()
}

func GetCoWorking(id int64) (cp *CoWorking, err error) {
	cp = &CoWorking{}
	cp.Product.Id = id

	o := orm.NewOrm()
	err = o.Read(&cp.Product)

	return
}

func GetAllCoWorking() (cps []*CoWorking, err error) {
	products, err := GetAllOfType(TYPE_CO_WORKING)
	if err != nil {
		return
	}

	cps = make([]*CoWorking, 0, len(products))
	for _, product := range products {
		cp := &CoWorking{
			Product: *product,
		}
		cps = append(cps, cp)
	}

	return
}

func GetAllCoWorkingAt(locationId int64) (cps []*CoWorking, err error) {
	products, err := GetAllOfTypeAt(locationId, TYPE_CO_WORKING)
	if err != nil {
		return
	}

	cps = make([]*CoWorking, 0, len(products))
	for _, product := range products {
		cp := &CoWorking{
			Product: *product,
		}
		cps = append(cps, cp)
	}

	return
}
