package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
)

type SpacePurchase struct {
	json.Marshaler
	json.Unmarshaler
	purchase Purchase
}

func (this *SpacePurchase) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.purchase)
}

func (this *SpacePurchase) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.purchase)
}

func CreateSpacePurchase(spacePurchase *SpacePurchase) (int64, error) {
	spacePurchase.purchase = Purchase{
		Created: time.Now(),
		Type:    PURCHASE_TYPE_SPACE_PURCHASE,
	}

	o := orm.NewOrm()
	return o.Insert(&spacePurchase.purchase)
}

func GetSpacePurchase(id int64) (spacePurchase *SpacePurchase, err error) {
	spacePurchase = &SpacePurchase{}
	spacePurchase.purchase.Id = id

	o := orm.NewOrm()
	err = o.Read(&spacePurchase.purchase)

	return
}

func GetAllSpacePurchases() (spacePurchases []*SpacePurchase, err error) {
	o := orm.NewOrm()
	r := new(SpacePurchase)
	var purchases []*Purchase
	_, err = o.QueryTable(r.purchase.TableName()).
		Filter("type", PURCHASE_TYPE_SPACE_PURCHASE).
		All(&purchases)
	if err != nil {
		return
	}
	spacePurchases = make([]*SpacePurchase, 0, len(purchases))
	for _, purchase := range purchases {
		spacePurchase := &SpacePurchase{
			purchase: *purchase,
		}
		spacePurchases = append(spacePurchases, spacePurchase)
	}
	return
}

func UpdateSpacePurchase(spacePurchase *SpacePurchase) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&spacePurchase.purchase)
	return
}

func DeleteSpacePurchase(id int64) (err error) {
	return DeletePurchase(id)
}
