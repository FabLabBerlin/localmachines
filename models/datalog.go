package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(DataLog))
}

type DataLog struct {
	Id           int64  `orm:"auto";"pk"`
	ChangedTable string `orm:"size(100)"`
	BeforeJson   string
	AfterJson    string
	Created      time.Time
}

func NewDataLog(changedTable string, before, after interface{}) (*DataLog, error) {
	var err error
	dataLog := DataLog{
		ChangedTable: changedTable,
		Created:      time.Now(),
	}
	var buf []byte
	buf, err = json.Marshal(before)
	if err != nil {
		return nil, err
	}
	dataLog.BeforeJson = string(buf)
	buf, err = json.Marshal(after)
	if err != nil {
		return nil, err
	}
	dataLog.AfterJson = string(buf)

	return &dataLog, nil
}

func CreateDataLog(dataLog *DataLog) (id int64, err error) {
	o := orm.NewOrm()
	return o.Insert(dataLog)
}

func (this *DataLog) Diff() (deltas []Delta, err error) {
	var beforeFields map[string]interface{}
	var afterFields map[string]interface{}

	buf := []byte(this.BeforeJson)
	err = json.Unmarshal(buf, &beforeFields)
	if err != nil {
		return
	}

	buf = []byte(this.AfterJson)
	err = json.Unmarshal(buf, &afterFields)
	if err != nil {
		return
	}

	if len(beforeFields) != len(afterFields) {
		return nil, fmt.Errorf("expecting len(beforeFields) == len(afterFields)")
	}

	deltas = make([]Delta, 0, len(beforeFields))
	for field, before := range beforeFields {
		var ok bool
		delta := Delta{
			FieldName: field,
			Before:    before,
		}
		delta.After, ok = afterFields[field]
		if !ok {
			return nil, fmt.Errorf("field %v not present in afterField", field)
		}
		isZero, err := delta.isZero()
		if err != nil {
			return nil, err
		}
		if !isZero {
			deltas = append(deltas, delta)
		}
	}

	return
}

func (this *DataLog) TableName() string {
	return "data_log"
}

type Delta struct {
	FieldName string
	Before    interface{}
	After     interface{}
}

func (this *Delta) isInsert() bool {
	return this.Before == nil && this.After != nil
}

func (this *Delta) isDelete() bool {
	return this.Before != nil && this.After == nil
}

func (this *Delta) isZero() (bool, error) {
	before, err := json.Marshal(this.Before)
	if err != nil {
		return false, err
	}
	after, err := json.Marshal(this.After)
	if err != nil {
		return false, err
	}
	return bytes.Equal(before, after), nil
}
