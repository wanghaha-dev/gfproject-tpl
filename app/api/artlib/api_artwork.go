package artlib

import (
	"myproject/utils"
	"strings"

	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"github.com/wanghaha-dev/gf/util/guid"
)

var Artwork = new(artwork)

type artwork struct{}

// List 获取数据列表
func (*artwork) List(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)

	whereMap := make(map[string]interface{})
	createName := r.GetString("create_name")
	createDate := r.GetString("create_date")

	if createName != "" {
		whereMap["create_name like ?"] = "%" + createName + "%"
	}
	if createDate != "" {
		whereMap["create_date like ?"] = "%" + createDate + "%"
	}

	total, err := g.DB("artlib").Model("al_work").Where(whereMap).Count()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	result, err := g.DB("artlib").Model("al_work").Where(whereMap).
		OrderDesc("create_date").
		Offset((page - 1) * limit).
		Limit(limit).
		FindAll()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	utils.AdminRespOkWithPageData(r, result, total)
}

// Add 新增数据
func (*artwork) Add(r *ghttp.Request) {
	params := r.GetMap()
	params["id"] = guid.S()

	rules, ruleErr := utils.GetAddFieldsRule("al_work")
	if ruleErr != nil {
		g.Log().Error(ruleErr)
		utils.RespFailWithMsg(r, ruleErr.Error())
		return
	}

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, insertErr := g.DB("artlib").Model("al_work").
		Insert(params)
	if insertErr != nil {
		g.Log().Error(insertErr)
		utils.RespFail(r)
	}

	utils.RespOk(r)
}

// DetailById 根据Id获取数据详细信息
func (*artwork) DetailById(r *ghttp.Request) {
	record, err := g.DB("artlib").Model("al_work").
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
func (*artwork) DeleteById(r *ghttp.Request) {
	_, err := g.DB("artlib").Model("al_work").
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
func (*artwork) UpdateById(r *ghttp.Request) {
	// 新数据
	params := r.GetMap()

	rules, ruleErr := utils.GetUpdateFieldsRule("al_work")
	if ruleErr != nil {
		g.Log().Error(ruleErr)
		utils.RespFailWithMsg(r, ruleErr.Error())
		return
	}

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, updateErr := g.DB("artlib").Model("al_work").
		Where("id=?", r.GetString("id")).
		Update(params)
	if updateErr != nil {
		g.Log().Error(updateErr)
		utils.RespFail(r)
		return
	}

	utils.RespOk(r)
}

func (*artwork) BatchUpdate(r *ghttp.Request) {
	data, err := r.GetJson()
	if err != nil {
		panic(err)
	}

	workIds := data.GetStrings("workIds")
	periodIds := data.GetStrings("periods")

	periodObjs, err := g.DB("artlib").Model("al_period").Where("id in (?)", periodIds).FindAll()
	utils.GetErrExit(r, err)

	workNameJoin := strings.Join(utils.ResultStrSlice(periodObjs, "name"), ",")
	workIdJoin := strings.Join(workIds, ",")
	_ = workIdJoin
	periodIdJoin := strings.Join(periodIds, ",")

	_, err = g.DB("artlib").Model("al_work").Where("id in (?)", workIds).Update(g.Map{
		"period_ids": periodIdJoin,
		"imp_period": workNameJoin,
	})
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*artwork) UpdateOne(r *ghttp.Request) {
	data, err := r.GetJson()
	utils.GetErrExit(r, err)

	workDesc := data.GetString("work_desc")
	workId := data.GetString("id")
	name := data.GetString("name")
	enName := data.GetString("en_name")
	authorName := data.GetString("imp_author_name")
	authorEnName := data.GetString("imp_author_enname")

	_, err = g.DB("artlib").Model("al_work").Where("id = ?", workId).Update(g.Map{
		"work_desc":         workDesc,
		"imp_author_name":   authorName,
		"imp_author_enname": authorEnName,
		"name":              name,
		"en_name":           enName,
	})
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

// Index 首页
func (*artwork) Index(r *ghttp.Request) {
	// todo
	utils.RespOkWithDataAndCache(r, "hello artwork")
}
