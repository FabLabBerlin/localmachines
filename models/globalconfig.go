package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
)

type GlobalConfig struct {
	Id    int64  `orm:"auto";"pk"`
	Name  string `orm:"size(100)"`
	Value string `orm:"type(text)"`
}

func (this *GlobalConfig) SetValueFloat(f float64) {
	this.Value = strconv.FormatFloat(f, 'f', -1, 64)
}

func (this *GlobalConfig) ValueFloat() (f float64, err error) {
	return strconv.ParseFloat(this.Value, 64)
}

func (this *GlobalConfig) SetValueInt(i int) {
	this.Value = strconv.Itoa(i)
}

func (this *GlobalConfig) ValueInt() (i int, err error) {
	return strconv.Atoi(this.Value)
}
