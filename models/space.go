package models

import (
	"github.com/astaxie/beego/orm"
)

type Space struct {
	Product Product
}

func CreateSpace(space *Space) (int64, error) {
	space.Product.Type = PURCHASE_TYPE_SPACE_RESERVATION
	return CreateProduct(&space.Product)
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
