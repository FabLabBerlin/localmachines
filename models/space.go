package models

import (
	"github.com/astaxie/beego/orm"
)

type Space struct {
	Product Product
}

func CreateSpace(name string) (space Space, err error) {
	space = Space{
		Product: Product{
			Name: name,
			Type: PRODUCT_TYPE_SPACE,
		},
	}
	id, err := CreateProduct(&space.Product)
	space.Product.Id = id
	return
}

func GetAllSpaces() (spaces []*Space, err error) {
	o := orm.NewOrm()
	var products []*Product
	_, err = o.QueryTable(new(Product).TableName()).
		Filter("type", PRODUCT_TYPE_SPACE).
		All(&products)

	if err != nil {
		return
	}

	spaces = make([]*Space, 0, len(products))
	for _, product := range products {
		space := &Space{
			Product: *product,
		}
		spaces = append(spaces, space)
	}

	return
}
