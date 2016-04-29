package products

import (
	"errors"
	"github.com/astaxie/beego/orm"
)

const TABLE_NAME = "products"

const (
	TYPE_TUTOR      = "tutor"
)

type Product struct {
	Id            int64
	LocationId    int64
	Type          string `orm:"size(100)"`
	Name          string `orm:"size(100)"`
	Price         float64
	PriceUnit     string `orm:"size(100)"`
	UserId        int64
	MachineSkills string `orm:"size(255)"` // JSON string [1, 3, 6]
	Comments      string `orm:"type(text)"`
	Archived      bool
}

type ProductList struct {
	Data []*Product
}

func (this *Product) TableName() string {
	return TABLE_NAME
}

func Create(product *Product) (int64, error) {
	if product.LocationId <= 0 {
		return 0, errors.New("LocationId must be > 0")
	}
	o := orm.NewOrm()
	return o.Insert(product)
}

func Get(productId int64) (product *Product, err error) {
	o := orm.NewOrm()
	product = &Product{Id: productId}
	err = o.Read(product)
	return
}

// Filter out archived products
func GetAll() (products []*Product, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).
		Exclude("archived", 1).
		All(&products)
	return
}

// Filter out archived products
func GetAllAt(locationId int64) (products []*Product, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).
		Filter("location_id", locationId).
		Exclude("archived", 1).
		All(&products)
	return
}

func GetAllOfType(productType string) (products []*Product, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).
		Filter("type", productType).
		Exclude("archived", 1).
		All(&products)
	return
}

func GetAllOfTypeAt(locationId int64, typ string) (ps []*Product, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(TABLE_NAME).
		Filter("location_id", locationId).
		Filter("type", typ).
		Exclude("archived", 1).
		All(&ps)
	return
}

func (product *Product) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(product)
	return
}

func (product *Product) Archive() (err error) {
	product.Archived = true
	return product.Update()
}

func init() {
	orm.RegisterModel(new(Product))
}
