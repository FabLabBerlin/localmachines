package products

import (
	"github.com/astaxie/beego/orm"
)

type CoWorkingProduct struct {
	Product Product
}

func CreateCoWorkingProduct(name string) (cp CoWorkingProduct, err error) {
	cp = CoWorkingProduct{
		Product: Product{
			Name: name,
			Type: PRODUCT_TYPE_CO_WORKING,
		},
	}
	id, err := CreateProduct(&cp.Product)
	cp.Product.Id = id
	return
}

func UpdateCoWorkingProduct(cp *CoWorkingProduct) (err error) {
	return UpdateProduct(&cp.Product)
}

func GetCoWorkingProduct(id int64) (cp *CoWorkingProduct, err error) {
	cp = &CoWorkingProduct{}
	cp.Product.Id = id

	o := orm.NewOrm()
	err = o.Read(&cp.Product)

	return
}

func GetAllCoWorkingProducts() (cps []*CoWorkingProduct, err error) {
	products, err := GetAllProductsOfType(PRODUCT_TYPE_CO_WORKING)
	if err != nil {
		return
	}

	cps = make([]*CoWorkingProduct, 0, len(products))
	for _, product := range products {
		cp := &CoWorkingProduct{
			Product: *product,
		}
		cps = append(cps, cp)
	}

	return
}
