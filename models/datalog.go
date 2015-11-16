package models

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
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

func DataLogSync() error {
	task := &DataSyncTask{}
	return task.run()
}

func GetAllDataLogs() (dataLogs []*DataLog, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("data_log").All(&dataLogs)
	return
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

func (this *Delta) Hash() (hash string, err error) {
	before, err := json.Marshal(this.Before)
	if err != nil {
		return
	}
	after, err := json.Marshal(this.After)
	if err != nil {
		return
	}
	data := bytes.Join([][]byte{before, after}, []byte{})
	hash = fmt.Sprintf("%x", sha1.Sum(data))
	return
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

type DataSyncTask struct {
	localLogs       []*DataLog
	localLogHashes  []string
	remoteLogs      map[string][]*DataLog
	remoteLogHashes []string
}

func (this *DataSyncTask) run() error {
	beego.Info("Running DataLogSync Task")
	this.fetchLocalLogs()
	switch beego.AppConfig.String("instancetype") {
	case "cloud":
		labs := strings.Split(beego.AppConfig.String("labs"), ",")
		if len(labs) == 0 {
			return fmt.Errorf("app.conf: no labs defined")
		}
		break
	case "lab":
		cloud := beego.AppConfig.String("cloud")
		if cloud == "" {
			return fmt.Errorf("app.conf: no cloud defined")
		}
		break
	default:
		return fmt.Errorf("app.conf: instancetype must be cloud or lab")
	}
	return nil
}

func (this *DataSyncTask) fetchLocalLogs() (err error) {
	this.localLogs, err = GetAllDataLogs()
	return
}
