package artlib

import (
	"myproject/utils"

	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"github.com/wanghaha-dev/gf/util/guid"
)

var AlPeriod = new(alPeriod)

type alPeriod struct{}

// List 获取数据列表
func (*alPeriod) List(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)

	total, err := g.DB("default").Model("al_period").Count()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	result, err := g.DB("default").Model("al_period").
		Offset((page - 1) * limit).
		Limit(limit).
		FindAll()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	utils.RespOkWithPageDataAndCache(r, result, total)
}

// All 获取所有数据不分页
func (*alPeriod) All(r *ghttp.Request) {
	result, err := g.DB("default").Model("al_period").
		FindAll()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	utils.RespOkWithDataAndCache(r, result)
}

// Add 新增数据
func (*alPeriod) Add(r *ghttp.Request) {
	params := r.GetMap()
	params["id"] = guid.S()

	rules, ruleErr := utils.GetAddFieldsRule("al_period")
	if ruleErr != nil {
		g.Log().Error(ruleErr)
		utils.RespFailWithMsg(r, ruleErr.Error())
		return
	}

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, insertErr := g.DB("default").Model("al_period").
		Insert(params)
	if insertErr != nil {
		g.Log().Error(insertErr)
		utils.RespFail(r)
	}

	utils.RespOk(r)
}

// DetailById 根据Id获取数据详细信息
func (*alPeriod) DetailById(r *ghttp.Request) {
	record, err := g.DB("default").Model("al_period").
		Where("id=?", r.GetString("id")).
		FindOne()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	utils.RespOkWithDataAndCache(r, record)
}

// DeleteById 根据Id删除数据
func (*alPeriod) DeleteById(r *ghttp.Request) {
	_, err := g.DB("default").Model("al_period").
		Where("id=?", r.GetString("id")).
		Delete()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	utils.RespOk(r)
}

// UpdateById 根据id修改数据
func (*alPeriod) UpdateById(r *ghttp.Request) {
	// 新数据
	params := r.GetMap()

	rules, ruleErr := utils.GetUpdateFieldsRule("al_period")
	if ruleErr != nil {
		g.Log().Error(ruleErr)
		utils.RespFailWithMsg(r, ruleErr.Error())
		return
	}

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, updateErr := g.DB("default").Model("al_period").
		Where("id=?", r.GetString("id")).
		Update(params)
	if updateErr != nil {
		g.Log().Error(updateErr)
		utils.RespFail(r)
		return
	}

	utils.RespOk(r)
}

// Index 首页
func (*alPeriod) Index(r *ghttp.Request) {
	// todo
	utils.RespOkWithDataAndCache(r, "hello alPeriod")
}
