package main

import (
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/memberships"
	"github.com/astaxie/beego/migration"
	"github.com/astaxie/beego/orm"
)

// DO NOT MODIFY
type MembershipsPopulateCategoryIds_20170112_143836 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MembershipsPopulateCategoryIds_20170112_143836{}
	m.Created = "20170112_143836"
	migration.Register("MembershipsPopulateCategoryIds_20170112_143836", m)
}

func handleMb(mb *memberships.Membership) (err error) {
	mIds, err := mb.AffectedMachineIds()
	if err != nil {
		return
	}
	cIds := make([]int64, 0, 10)
	for _, mId := range mIds {
		m, err := machine.Get(mId)
		if err != nil {
			return err
		}
		cIds = append(cIds, m.TypeId)
	}
	cIds = uniq(cIds)
	if err = mb.SetAffectedCategoryIds(cIds); err != nil {
		return
	}

	if err = mb.Update(); err != nil {
		return
	}

	return
}

func uniq(ids []int64) (res []int64) {
	h := make(map[int64]bool)

	for _, id := range ids {
		h[id] = true
	}

	res = make([]int64, 0, len(res))
	for id := range h {
		res = append(res, id)
	}

	return
}

// Run the migrations
func (m *MembershipsPopulateCategoryIds_20170112_143836) Up() {
	o := orm.NewOrm()
	o.Begin()

	var mbs []*memberships.Membership
	_, err := o.QueryTable("membership").
		All(&mbs)
	if err != nil {
		panic(err.Error())
	}

	for _, mb := range mbs {
		if err := handleMb(mb); err != nil {
			o.Rollback()
			panic(err.Error())
		}
	}

	if err := o.Commit(); err != nil {
		panic(err.Error())
	}
}

// Reverse the migrations
func (m *MembershipsPopulateCategoryIds_20170112_143836) Down() {
	// no necessity to reverse anything
}
