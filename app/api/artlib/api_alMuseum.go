package artlib

import (
	"myproject/utils"

	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"github.com/wanghaha-dev/gf/util/guid"
)

var Almuseum = new(alMuseum)

type alMuseum struct{}

// List 获取数据列表
func (*alMuseum) List(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)

	total, err := g.DB("artlib").Model("al_institute").Count()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	result, err := g.DB("artlib").Model("al_institute").
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
func (*alMuseum) All(r *ghttp.Request) {
	result, err := g.DB("artlib").Model("al_institute").
		FindAll()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	utils.RespOkWithDataAndCache(r, result)
}

// Add 新增数据
func (*alMuseum) Add(r *ghttp.Request) {
	params := r.GetMap()
	params["id"] = guid.S()

	rules, ruleErr := utils.GetAddFieldsRule("al_institute")
	if ruleErr != nil {
		g.Log().Error(ruleErr)
		utils.RespFailWithMsg(r, ruleErr.Error())
		return
	}

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, insertErr := g.DB("artlib").Model("al_institute").
		Insert(params)
	if insertErr != nil {
		g.Log().Error(insertErr)
		utils.RespFail(r)
	}

	utils.RespOk(r)
}

// DetailById 根据Id获取数据详细信息
func (*alMuseum) DetailById(r *ghttp.Request) {
	record, err := g.DB("artlib").Model("al_institute").
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
func (*alMuseum) DeleteById(r *ghttp.Request) {
	_, err := g.DB("artlib").Model("al_institute").
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
func (*alMuseum) UpdateById(r *ghttp.Request) {
	// 新数据
	params := r.GetMap()

	rules, ruleErr := utils.GetUpdateFieldsRule("al_institute")
	if ruleErr != nil {
		g.Log().Error(ruleErr)
		utils.RespFailWithMsg(r, ruleErr.Error())
		return
	}

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, updateErr := g.DB("artlib").Model("al_institute").
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
func (*alMuseum) Index(r *ghttp.Request) {
	// todo
	utils.RespOkWithDataAndCache(r, "hello alMuseum")
}
