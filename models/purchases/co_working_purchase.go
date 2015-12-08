package purchases

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
)

type CoWorkingPurchase struct {
	json.Marshaler
	json.Unmarshaler
	purchase Purchase
}

func (this *CoWorkingPurchase) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.purchase)
}

func (this *CoWorkingPurchase) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.purchase)
}

func CreateCoWorkingPurchase(cp *CoWorkingPurchase) (int64, error) {
	cp.purchase = Purchase{
		Created:   time.Now(),
		Type:      PURCHASE_TYPE_CO_WORKING,
		TimeStart: time.Now(),
		TimeEnd:   time.Now(),
		PriceUnit: "hour",
	}

	o := orm.NewOrm()
	return o.Insert(&cp.purchase)
}

func GetCoWorkingPurchase(id int64) (cp *CoWorkingPurchase, err error) {
	cp = &CoWorkingPurchase{}
	cp.purchase.Id = id

	o := orm.NewOrm()
	err = o.Read(&cp.purchase)

	return
}

func GetAllCoWorkingPurchases() (cps []*CoWorkingPurchase, err error) {
	purchases, err := GetAllPurchasesOfType(PURCHASE_TYPE_CO_WORKING)
	if err != nil {
		return
	}
	cps = make([]*CoWorkingPurchase, 0, len(purchases))
	for _, purchase := range purchases {
		cp := &CoWorkingPurchase{
			purchase: *purchase,
		}
		cps = append(cps, cp)
	}
	return
}

func UpdateCoWorkingPurchase(cp *CoWorkingPurchase) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&cp.purchase)
	return
}

func DeleteCoWorkingPurchase(id int64) (err error) {
	return DeletePurchase(id)
}
