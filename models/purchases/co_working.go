package purchases

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
)

type CoWorking struct {
	json.Marshaler
	json.Unmarshaler
	purchase Purchase
}

func (this *CoWorking) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.purchase)
}

func (this *CoWorking) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.purchase)
}

func CreateCoWorking(cp *CoWorking) (int64, error) {
	cp.purchase = Purchase{
		Created:   time.Now(),
		Type:      TYPE_CO_WORKING,
		TimeStart: time.Now(),
		TimeEnd:   time.Now(),
		PriceUnit: "hour",
	}

	o := orm.NewOrm()
	return o.Insert(&cp.purchase)
}

func GetCoWorking(id int64) (cp *CoWorking, err error) {
	cp = &CoWorking{}
	cp.purchase.Id = id

	o := orm.NewOrm()
	err = o.Read(&cp.purchase)

	return
}

func GetAllCoWorking() (cps []*CoWorking, err error) {
	purchases, err := GetAllOfType(TYPE_CO_WORKING)
	if err != nil {
		return
	}
	cps = make([]*CoWorking, 0, len(purchases))
	for _, purchase := range purchases {
		cp := &CoWorking{
			purchase: *purchase,
		}
		cps = append(cps, cp)
	}
	return
}

func UpdateCoWorking(cp *CoWorking) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&cp.purchase)
	return
}

func DeleteCoWorking(id int64) (err error) {
	return Delete(id)
}
