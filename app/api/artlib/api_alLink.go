package artlib

import (
	"myproject/utils"

	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"github.com/wanghaha-dev/gf/util/guid"
)

var AlLink = new(alLink)

type alLink struct{}

// List 获取数据列表
func (*alLink) List(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)

	total, err := g.DB("artlib").Model("al_link").Count()
	utils.GetErrExit(r, err)

	result, err := g.DB("artlib").Model("al_link").
		Offset((page - 1) * limit).
		Limit(limit).
		FindAll()
	utils.GetErrExit(r, err)

	utils.RespOkWithPageDataAndCache(r, result, total)
}

// Add 新增数据
func (*alLink) Add(r *ghttp.Request) {
	params := r.GetMap()
	params["id"] = guid.S()

	rules, ruleErr := utils.GetAddFieldsRule("al_link")
	utils.GetErrWithMsgExit(r, ruleErr, ruleErr.Error())

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, insertErr := g.DB("artlib").Model("al_link").
		Insert(params)
	if insertErr != nil {
		g.Log().Error(insertErr)
		utils.RespFail(r)
	}

	utils.RespOk(r)
}

// DetailById 根据Id获取数据详细信息
func (*alLink) DetailById(r *ghttp.Request) {
	record, err := g.DB("artlib").Model("al_link").
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
func (*alLink) DeleteById(r *ghttp.Request) {
	_, err := g.DB("artlib").Model("al_link").
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
func (*alLink) UpdateById(r *ghttp.Request) {
	// 新数据
	params := r.GetMap()

	rules, ruleErr := utils.GetUpdateFieldsRule("al_link")
	if ruleErr != nil {
		g.Log().Error(ruleErr)
		utils.RespFailWithMsg(r, ruleErr.Error())
		return
	}

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, updateErr := g.DB("artlib").Model("al_link").
		Where("id=?", r.GetString("id")).
		Update(params)
	utils.GetErrExit(r, updateErr)

	utils.RespOk(r)
}

// Index 首页
func (*alLink) Index(r *ghttp.Request) {
	// todo
	utils.RespOkWithDataAndCache(r, "hello alLink")
}
