package artlib

import (
	"github.com/wanghaha-dev/gf/container/gvar"
	"github.com/wanghaha-dev/gf/database/gdb"
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"github.com/wanghaha-dev/gf/util/guid"
	"myproject/utils"
)

var Alwork = new(alWork)

type alWork struct{}

// List 获取数据列表
func (*alWork) List(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)
	subjectCode := r.GetString("subjectCode")
	artistId := r.GetString("artistId")

	whereMap := make(map[string]interface{})

	if subjectCode == "-1" || subjectCode == "undefined" {

	} else if subjectCode != "" {
		whereMap["subject_codes like ?"] = "%" + subjectCode + "%"
	}

	if artistId != "" && artistId != "undefined" {
		whereMap["artist_ids like ?"] = "%" + artistId + "%"
	}

	total, err := g.DB("artlib").Model("al_work").
		Where("publish_flag=?", "1").
		Where(whereMap).Count()
	utils.GetErrExit(r, err)

	result, err := g.DB("artlib").Model("al_work").
		Where("publish_flag=?", "1").
		Where(whereMap).
		OrderDesc("create_date").
		Offset((page - 1) * limit).
		Limit(limit).
		FindAll()
	utils.GetErrExit(r, err)

	utils.RespOkWithPageDataAndCache(r, result, total)
}

// All 获取所有数据不分页
func (*alWork) All(r *ghttp.Request) {
	result, err := g.DB("artlib").Model("al_work").
		Where("publish_flag=?", "1").
		FindAll()
	utils.GetErrExit(r, err)

	utils.RespOkWithDataAndCache(r, result)
}

// Add 新增数据
func (*alWork) Add(r *ghttp.Request) {
	params := r.GetMap()
	params["id"] = guid.S()

	rules, ruleErr := utils.GetAddFieldsRule("al_work")
	utils.GetErrWithMsgExit(r, ruleErr, ruleErr.Error())

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, insertErr := g.DB("artlib").Model("al_work").
		Insert(params)
	utils.GetErrExit(r, insertErr)

	utils.RespOk(r)
}

// DetailById 根据Id获取数据详细信息
func (*alWork) DetailById(r *ghttp.Request) {
	record, err := g.DB("artlib").Model("al_work").
		Where("id=?", r.GetString("id")).
		FindOne()
	utils.GetErrExit(r, err)

	utils.RespOkWithData(r, record)
}

// DeleteById 根据Id删除数据
func (*alWork) DeleteById(r *ghttp.Request) {
	_, err := g.DB("artlib").Model("al_work").
		Where("id=?", r.GetString("id")).
		Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

// UpdateById 根据id修改数据
func (*alWork) UpdateById(r *ghttp.Request) {
	// 新数据
	params := r.GetMap()

	rules, ruleErr := utils.GetUpdateFieldsRule("al_work")
	utils.GetErrWithMsgExit(r, ruleErr, ruleErr.Error())

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, updateErr := g.DB("artlib").Model("al_work").
		Where("id=?", r.GetString("id")).
		Update(params)
	utils.GetErrExit(r, updateErr)

	utils.RespOk(r)
}

// Index 首页
func (*alWork) Index(r *ghttp.Request) {
	// todo
	utils.RespOkWithData(r, "hello alWork")
}

// Download 根据Id获取数据列表
func (*alWork) Download(r *ghttp.Request) {
	records, err := g.DB("artlib").Model("al_work_files").
		Where("work_id=?", r.GetString("id")).
		FindAll()
	utils.GetErrExit(r, err)

	utils.RespOkWithData(r, records)
}

// FirstByCode 根据code获取数据详细信息
func (*alWork) FirstByCode(r *ghttp.Request) {
	subjectCode := r.GetString("subjectCode")
	whereMap := make(map[string]interface{})

	if subjectCode == "-1" {

	} else if subjectCode != "" {
		whereMap["subject_codes like ?"] = "%" + subjectCode + "%"
	}

	work, err := g.DB("artlib").Model("al_work").
		Where("publish_flag=?", "1").
		Where(whereMap).
		FindOne()
	utils.GetErrExit(r, err)

	record, err := g.DB("artlib").Model("al_work_files").
		Where("work_id =?", work["id"]).
		FindOne()
	utils.GetErrExit(r, err)

	utils.RespOkWithDataAndCache(r, record)
}

// TypesCount 获取分类下的艺术品数量
func (*alWork) TypesCount(r *ghttp.Request) {
	// 拿到所有分类
	typeData, err := g.Model("z_mini_type_codes").Fields("code, name").Where("parent_code=?", "02").FindAll()
	utils.GetErrExit(r, err)

	typeDataNotEmpty := make(gdb.Result, 0)

	// 遍历分类, 获取每个分类的作品数量
	// 并且获取每个分类的代表图片
	for _, item := range typeData {
		count, err := g.Model("al_work").
			Where("publish_flag=?", "1").
			Where("subject_codes like ?", "%"+item["code"].String()+"%").
			Count()
		utils.GetErrExit(r, err)

		item["count"] = gvar.New(count)

		// 如果分类下有作品的话, 就拿第一个作品
		if count != 0 {
			firstWork, err := g.Model("al_work").
				Where("publish_flag=?", "1").
				Where("subject_codes like ?", item["code"]).FindOne()
			utils.GetErrExit(r, err)

			item["image"] = gvar.New(firstWork["imgs"])
			typeDataNotEmpty = append(typeDataNotEmpty, item)
		}
	}

	// 拿到总数
	total, err := g.Model("al_work").
		Where("publish_flag=?", "1").
		Count()
	utils.GetErrExit(r, err)

	// 若总作品数为0则不需要加了
	if total != 0 {
		first, err := g.Model("al_work").
			Where("publish_flag=?", "1").
			FindOne()
		utils.GetErrExit(r, err)

		allData := make(gdb.Record)
		allData["count"] = gvar.New(total)
		allData["name"] = gvar.New("全部")
		allData["image"] = gvar.New(first["imgs"])
		typeDataNotEmpty = append(typeDataNotEmpty, allData)
	}

	// add cache
	utils.RespOkWithDataAndCache(r, typeDataNotEmpty)
}
