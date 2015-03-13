package models

import (
	"errors"
	"strconv"
	"time"
)

var (
	TestObjects map[string]*TestObject
)

type TestObject struct {
	ObjectId   string
	Score      int64
	PlayerName string
}

func init() {
	TestObjects = make(map[string]*TestObject)
	TestObjects["hjkhsbnmn123"] = &TestObject{"hjkhsbnmn123", 100, "astaxie"}
	TestObjects["mjjkxsxsaa23"] = &TestObject{"mjjkxsxsaa23", 101, "someone"}
}

func AddOne(object TestObject) (ObjectId string) {
	object.ObjectId = "astaxie" + strconv.FormatInt(time.Now().UnixNano(), 10)
	TestObjects[object.ObjectId] = &object
	return object.ObjectId
}

func GetOne(ObjectId string) (object *TestObject, err error) {
	if v, ok := TestObjects[ObjectId]; ok {
		return v, nil
	}
	return nil, errors.New("ObjectId Not Exist")
}

func GetAll() map[string]*TestObject {
	return TestObjects
}

func Update(ObjectId string, Score int64) (err error) {
	if v, ok := TestObjects[ObjectId]; ok {
		v.Score = Score
		return nil
	}
	return errors.New("ObjectId Not Exist")
}

func Delete(ObjectId string) {
	delete(TestObjects, ObjectId)
}

