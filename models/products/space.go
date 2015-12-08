package products

import (
	"github.com/astaxie/beego"
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

func UpdateSpace(space *Space) (err error) {
	return UpdateProduct(&space.Product)
}

func GetSpace(id int64) (space *Space, err error) {
	space = &Space{}
	space.Product.Id = id

	beego.Info("GetSpace: id=", id)

	o := orm.NewOrm()
	err = o.Read(&space.Product)

	return
}

func GetAllSpaces() (spaces []*Space, err error) {
	products, err := GetAllProductsOfType(PRODUCT_TYPE_SPACE)
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
