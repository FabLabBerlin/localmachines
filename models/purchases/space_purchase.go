package purchases

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
)

type Space struct {
	json.Marshaler
	json.Unmarshaler
	Purchase
}

func NewSpace(locationId int64) *Space {
	return &Space{
		Purchase: Purchase{
			LocationId: locationId,
			Created:    time.Now(),
			Type:       TYPE_SPACE,
			TimeStart:  time.Now(),
			TimeEnd:    time.Now(),
			PriceUnit:  "hour",
		},
	}
}

func (this *Space) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.Purchase)
}

func (this *Space) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &this.Purchase)
}

func (this *Space) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&this.Purchase)
	return
}

func (this *Space) Save() (err error) {
	o := orm.NewOrm()
	if this.Id == 0 {
		_, err = o.Insert(&this.Purchase)
	} else {
		err = this.Update()
	}
	return
}

func GetSpace(id int64) (spacePurchase *Space, err error) {
	spacePurchase = &Space{}
	spacePurchase.Id = id

	o := orm.NewOrm()
	err = o.Read(&spacePurchase.Purchase)

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
			Purchase: *purchase,
		}
		spacePurchases = append(spacePurchases, spacePurchase)
	}
	return
}

func GetAllSpaceAt(locationId int64) (spacePurchases []*Space, err error) {
	purchases, err := GetAllOfTypeAt(locationId, TYPE_SPACE)
	if err != nil {
		return
	}
	spacePurchases = make([]*Space, 0, len(purchases))
	for _, purchase := range purchases {
		spacePurchase := &Space{
			Purchase: *purchase,
		}
		spacePurchases = append(spacePurchases, spacePurchase)
	}
	return
}

func DeleteSpace(id int64) (err error) {
	return Delete(id)
}
