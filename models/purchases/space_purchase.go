package purchases

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
)

type Space struct {
	json.Marshaler
	json.Unmarshaler
	purchase Purchase
}

func (this *Space) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.purchase)
}

func (this *Space) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.purchase)
}

func CreateSpace(spacePurchase *Space) (int64, error) {
	spacePurchase.purchase = Purchase{
		Created:   time.Now(),
		Type:      TYPE_SPACE,
		TimeStart: time.Now(),
		TimeEnd:   time.Now(),
		PriceUnit: "hour",
	}

	o := orm.NewOrm()
	return o.Insert(&spacePurchase.purchase)
}

func GetSpace(id int64) (spacePurchase *Space, err error) {
	spacePurchase = &Space{}
	spacePurchase.purchase.Id = id

	o := orm.NewOrm()
	err = o.Read(&spacePurchase.purchase)

	return
}

func GetAllSpace() (spacePurchases []*Space, err error) {
	purchases, err := GetAllOfType(TYPE_SPACE)
	if err != nil {
		return
	}
	spacePurchases = make([]*Space, 0, len(purchases))
	for _, purchase := range purchases {
		spacePurchase := &Space{
			purchase: *purchase,
		}
		spacePurchases = append(spacePurchases, spacePurchase)
	}
	return
}

func UpdateSpace(spacePurchase *Space) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&spacePurchase.purchase)
	return
}

func DeleteSpace(id int64) (err error) {
	return Delete(id)
}
