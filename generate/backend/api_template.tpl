package {{.PackageName}}

import (
	"myproject/utils"
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"github.com/wanghaha-dev/gf/util/guid"
)

var {{.StructNameUpper}} = new({{.StructName}})

type {{.StructName}} struct{}

// List 获取数据列表
func (*{{.StructName}}) List(r *ghttp.Request) {
	/*
		Condition Query
	*/
	whereMap := make(map[string]interface{})
	params1 := r.GetString("params1")

	if params1 != "" {
	// add whereMap
	// whereMap["params1"] = params1
	}

    // validator params
    /*
        rules := make(map[string]string)
        rules["params1"] = "required"
        checkErr := g.Validator().Rules(rules).CheckMap(whereMap)
        if checkErr != nil {
            utils.RespFailWithMsg(r, checkErr.FirstString())
        }
    */

	page, limit := utils.GetPageArgs(r)

	total, err := g.DB("{{.DBName}}").Model("{{.TableName}}").Where(whereMap).Count()
	utils.GetErrExit(r, err)

	result, err := g.DB("{{.DBName}}").Model("{{.TableName}}").
	Where(whereMap).
	Offset((page - 1) * limit).
	Limit(limit).
	FindAll()
	utils.GetErrExit(r, err)

    {{if .Cache}}
    utils.RespOkWithPageDataAndCache(r, result, total)
    {{else}}
    utils.RespOkWithPageData(r, result, total)
    {{end}}
}

// All 获取所有数据不分页
func (*{{.StructName}}) All(r *ghttp.Request) {
	/*
		Condition Query
	*/
	whereMap := make(map[string]interface{})
	params1 := r.GetString("params1")

	if params1 != "" {
	// add whereMap
	// whereMap["params1"] = params1
	}

	// validator params
    /*
        rules := make(map[string]string)
        rules["params1"] = "required"
        checkErr := g.Validator().Rules(rules).CheckMap(whereMap)
        if checkErr != nil {
            utils.RespFailWithMsg(r, checkErr.FirstString())
        }
    */

	result, err := g.DB("{{.DBName}}").Model("{{.TableName}}").
	Where(whereMap).
	FindAll()
	utils.GetErrExit(r, err)

	{{if .Cache}}
    utils.RespOkWithDataAndCache(r, result)
    {{else}}
    utils.RespOkWithData(r, result)
    {{end}}
}

// Add 新增数据
func (*{{.StructName}}) Add(r *ghttp.Request) {
	params := r.GetMap()
	params["id"] = guid.S()

	rules, ruleErr := utils.GetAddFieldsRule("{{.TableName}}")
	utils.GetErrWithMsgExit(r, ruleErr, ruleErr.Error())

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, insertErr := g.DB("{{.DBName}}").Model("{{.TableName}}").
	Insert(params)
	utils.GetErrExit(r, insertErr)

	utils.RespOk(r)
}

// DetailById 根据Id获取数据详细信息
func (*{{.StructName}}) DetailById(r *ghttp.Request) {
	/*
		Condition Query
	*/
	whereMap := make(map[string]interface{})
	id := r.GetString("id")

	if id != "" {
	// add whereMap
		whereMap["id"] = id
	}

	// validator params
    /*
        rules := make(map[string]string)
        rules["params1"] = "required"
        checkErr := g.Validator().Rules(rules).CheckMap(whereMap)
        if checkErr != nil {
            utils.RespFailWithMsg(r, checkErr.FirstString())
        }
    */

	record, err := g.DB("{{.DBName}}").Model("{{.TableName}}").
	Where(whereMap).
	FindOne()
	utils.GetErrExit(r, err)


	{{if .Cache}}
    utils.RespOkWithDataAndCache(r, record)
    {{else}}
    utils.RespOkWithData(r, record)
    {{end}}
}

// DeleteById 根据Id删除数据
func (*{{.StructName}}) DeleteById(r *ghttp.Request) {
	/*
		Condition Query
	*/
	whereMap := make(map[string]interface{})
	id := r.GetString("id")

	if id != "" {
		// add whereMap
		whereMap["id"] = id
	}

	// validator params
    /*
        rules := make(map[string]string)
        rules["params1"] = "required"
        checkErr := g.Validator().Rules(rules).CheckMap(whereMap)
        if checkErr != nil {
            utils.RespFailWithMsg(r, checkErr.FirstString())
        }
    */

	_, err := g.DB("{{.DBName}}").Model("{{.TableName}}").
	Where(whereMap).
	Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

// UpdateById 根据id修改数据
func (*{{.StructName}}) UpdateById(r *ghttp.Request) {
	// 新数据
	params := r.GetMap()

	rules, ruleErr := utils.GetUpdateFieldsRule("{{.TableName}}")
	utils.GetErrWithMsgExit(r, ruleErr, ruleErr.Error())

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, updateErr := g.DB("{{.DBName}}").Model("{{.TableName}}").
	Where("id=?", r.GetString("id")).
	Update(params)
	utils.GetErrExit(r, updateErr)

	utils.RespOk(r)
}

// Index 首页
func (*{{.StructName}}) Index(r *ghttp.Request) {
	// todo
	utils.RespOkWithData(r, "hello {{.StructName}}")
}
