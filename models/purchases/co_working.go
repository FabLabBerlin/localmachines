package purchases

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
)

type CoWorking struct {
	json.Marshaler
	json.Unmarshaler
	Purchase
}

func (this *CoWorking) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.Purchase)
}

func (this *CoWorking) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.Purchase)
}

func CreateCoWorking(cp *CoWorking) (int64, error) {
	cp.Created = time.Now()
	cp.Type = TYPE_CO_WORKING
	cp.TimeStart = time.Now()
	cp.TimeEnd = time.Now()
	cp.PriceUnit = "hour"

	o := orm.NewOrm()
	return o.Insert(&cp.Purchase)
}

func GetCoWorking(id int64) (cp *CoWorking, err error) {
	cp = &CoWorking{}
	cp.Purchase.Id = id

	o := orm.NewOrm()
	err = o.Read(&cp.Purchase)

	return
}

func GetAllCoWorking() (cps []*CoWorking, err error) {
	purchases, err := GetAllOfType(TYPE_CO_WORKING)
	if err != nil {
		return
	}
	cps = make([]*CoWorking, 0, len(purchases))
	for _, Purchase := range purchases {
		cp := &CoWorking{
			Purchase: *Purchase,
		}
		cps = append(cps, cp)
	}
	return
}

func (cp *CoWorking) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&cp.Purchase)
	return
}

func DeleteCoWorking(id int64) (err error) {
	return Delete(id)
}
