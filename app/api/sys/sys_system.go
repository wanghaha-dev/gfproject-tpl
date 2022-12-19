package sys

import (
	"encoding/json"
	"fmt"
	"myproject/utils"
	"myproject/utils/captcha"
	"myproject/utils/gjwt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wanghaha-dev/gf/container/gvar"
	"github.com/wanghaha-dev/gf/database/gdb"
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"github.com/wanghaha-dev/gf/util/gconv"
)

var Sys = new(sys)

type sys struct{}

// Index 首页
func (*sys) Index(r *ghttp.Request) {
	// todo
	utils.RespOkWithData(r, "hello user")
}

// Login 登录
func (*sys) Login(r *ghttp.Request) {
	// 新数据
	params := r.GetMap()

	rules := make(map[string]interface{})
	rules["username"] = "required"
	rules["password"] = "required"

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	obj, err := g.Model("sys_user").Where(g.Map{
		"username": params["username"],
		"password": params["password"],
	}).FindOne()

	if err != nil {
		utils.RespFail(r)
		g.Log().Error(err)
		return
	}

	if obj.IsEmpty() {
		utils.RespFailWithMsg(r, "用户名或密码错误")
		return
	}

	token, err := gjwt.GenerateToken(jwt.MapClaims{
		"userId":   obj["userId"],
		"username": obj["username"],
	})
	if err != nil {
		utils.RespFail(r)
		g.Log().Error(err)
	}

	utils.RespOkWithData(r, g.Map{
		"access_token": token,
	})
}

func (*sys) Captcha(r *ghttp.Request) {
	code, data := captcha.GetBase64Encode(4)
	utils.RespOkWithData(r, g.Map{
		"text":   code,
		"base64": data,
	})
}

func (*sys) AuthUser(r *ghttp.Request) {
	// parse token get data
	tokenData, err := gjwt.ParseToken(r.GetHeader("Authorization"))
	fmt.Println("tokenData:", tokenData)
	utils.GetErrWithMsgExit(r, err, "parse token error.")

	// get roles for user
	userObj, err := g.Model("sys_user").Where("userId=?", tokenData["userId"]).FindOne()
	utils.GetErrExit(r, err)

	if userObj.IsEmpty() {
		utils.RespFailWith401Msg(r, "data is not exists.")
		g.Log().Error(err)
		return
	}

	// get menus for roles
	roleIds := strings.Split(userObj["roles"].String(), ",")
	roleMenus, err := g.Model("sys_role_menu").Where("roleId IN(?)", roleIds).FindAll()
	utils.GetErrExit(r, err)

	menuIds := utils.ResultStrSlice(roleMenus, "menuIds")

	var menus gdb.Result

	// 如果用户有超级管理员角色的话直接给所有菜单，所有权限放行
	if utils.HasString("9999", roleIds) {
		fmt.Println("您是超级管理员")
		// get menuInfo for menus
		menus, err = g.Model("sys_menu").FindAll()
		utils.GetErrExit(r, err)
	} else {
		// get menuInfo for menus
		menus, err = g.Model("sys_menu").Where("menuId in (?)", menuIds).FindAll()
		utils.GetErrExit(r, err)
	}

	// get roles
	roles, err := g.Model("sys_role").Where("roleId in (?)", roleIds).FindAll()
	utils.GetErrExit(r, err)

	// user roles
	if roles.IsEmpty() {
		userObj["roles"] = gvar.New([]g.Map{})
	} else {
		userObj["roles"] = gvar.New(roles)
	}

	// user authorities
	if roles.IsEmpty() {
		userObj["authorities"] = gvar.New([]g.Map{})
	} else {
		userObj["authorities"] = gvar.New(menus)
	}

	utils.RespOkWithData(r, userObj)
}

func (*sys) UserPageList(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)

	whereMap := make(map[string]interface{})
	if r.GetString("organizationId") != "" {
		whereMap["organizationId"] = r.GetString("organizationId")
	}

	users, err := g.Model("sys_user").Where(whereMap).Offset((page - 1) * limit).Limit(limit).FindAll()
	utils.GetErrExit(r, err)

	total, err := g.Model("sys_user").Where(whereMap).Count()
	utils.GetErrExit(r, err)

	for _, userObj := range users {
		roleIds := strings.Split(userObj["roles"].String(), ",")

		// get roles
		roles, err := g.Model("sys_role").Where("roleId in (?)", roleIds).FindAll()
		utils.GetErrExit(r, err)

		for _, roleObj := range roles {
			roleObj["userId"] = userObj["userId"]
		}

		// start 找打字典里的 sex 字段拿到里面的值
		dictObj, err := g.Model("sys_dictionary").Where("dictCode=?", "sex").FindOne()
		utils.GetErrExit(r, err)

		if !dictObj.IsEmpty() {
			dictDataObj, err := g.Model("sys_dictionary_data").
				Where("dictId=?", dictObj["dictId"]).
				Where("dictDataCode=?", userObj["sex"]).FindOne()
			utils.GetErrExit(r, err)
			userObj["sexName"] = dictDataObj["dictDataName"]
		}
		// end

		userObj["roles"] = gvar.New(roles)
	}

	if users.IsEmpty() {
		utils.RespOkWithData(r, []g.Map{})
	}

	utils.RespOkWithPageData(r, users, total)
}

func (*sys) RoleList(r *ghttp.Request) {
	// get roles
	roles, err := g.Model("sys_role").FindAll()
	utils.GetErrExit(r, err)

	utils.RespOkWithData(r, roles)
}

func (*sys) RolePageList(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)

	total, err := g.Model("sys_role").Count()
	utils.GetErrExit(r, err)

	// get roles
	roles, err := g.Model("sys_role").Offset((page - 1) * limit).Limit(limit).FindAll()
	utils.GetErrExit(r, err)

	utils.RespOkWithPageData(r, roles, total)
}

func (*sys) AddUser(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"username": "required",
		"nickname": "required",
		"password": "required",
		"sex":      "required",
		"roleIds":  "required",
	}

	var roleIdList []string
	for _, item := range data["roleIds"].([]interface{}) {
		roleId := item.(json.Number).String()
		roleIdList = append(roleIdList, roleId)
	}
	roles := strings.Join(roleIdList, ",")
	data["roles"] = roles

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	if data["birthday"] == "" {
		data["birthday"] = nil
	}

	// 如果机构为空默认为1
	if data["organizationId"] == "" {
		data["organizationId"] = 1
	}

	_, err := g.Model("sys_user").Insert(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) DeleteUser(r *ghttp.Request) {
	if r.GetString("userId") == "9999" {
		utils.RespFailWithMsg(r, "管理员用户不可删除")
		return
	}

	_, err := g.Model("sys_user").Where("userId=?", r.GetString("userId")).Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) BatchDeleteUser(r *ghttp.Request) {
	getJson, err := r.GetJson()
	utils.GetErrExit(r, err)
	userIdInters := getJson.Array()
	userIds := utils.IS2SS(userIdInters)

	if utils.HasString("9999", userIds) {
		utils.RespFailWithMsg(r, "管理员用户不能被删除！")
		return
	}

	_, err = g.Model("sys_user").Where("userId in (?)", userIds).Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) UpdateUser(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"userId":   "required",
		"username": "required",
		"nickname": "required",
		"sex":      "required",
		"roleIds":  "required",
	}

	if gconv.String(data["userId"]) == "9999" {
		utils.RespFailWithMsg(r, "管理员用户不可修改")
		return
	}

	var roleIdList []string
	for _, item := range data["roleIds"].([]interface{}) {
		roleId := item.(json.Number).String()
		roleIdList = append(roleIdList, roleId)
	}
	roles := strings.Join(roleIdList, ",")
	data["roles"] = roles

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_user").Replace(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) UpdateRole(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"roleId":   "required",
		"roleCode": "required",
		"roleName": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	roleObj, err := g.Model("sys_role").Where("roleId=?", r.GetString("roleId")).FindOne()
	utils.GetErrExit(r, err)

	if !roleObj.IsEmpty() {
		if roleObj["roleCode"].String() == "administrator" {
			utils.RespFailWithMsg(r, "超级管理员不能被修改")
			return
		}
	}

	_, err = g.Model("sys_role").Replace(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) AddRole(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"roleCode": "required",
		"roleName": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_role").Insert(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) DeleteRole(r *ghttp.Request) {
	roleObj, err := g.Model("sys_role").Where("roleId=?", r.GetString("roleId")).FindOne()
	utils.GetErrExit(r, err)

	if !roleObj.IsEmpty() {
		if roleObj["roleCode"].String() == "administrator" {
			utils.RespFailWithMsg(r, "超级管理员不能被删除")
			return
		}
	}

	_, err = g.Model("sys_role").Where("roleId=?", r.GetString("roleId")).Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) BatchDeleteRole(r *ghttp.Request) {
	getJson, err := r.GetJson()
	utils.GetErrExit(r, err)
	roleIdInters := getJson.Array()
	roleIds := utils.IS2SS(roleIdInters)

	if utils.HasString("9999", roleIds) {
		utils.RespFailWithMsg(r, "超级管理员角色不能被删除！")
		return
	}

	_, err = g.Model("sys_role").Where("roleId in (?)", roleIds).Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) Existence(r *ghttp.Request) {
	field := r.GetString("field")
	value := r.GetString("value")
	id := r.GetString("id")
	_ = id

	_, err := g.Model("sys_user").Where(fmt.Sprintf("%s=?", field), value).FindOne()
	utils.GetErrExit(r, err)

	utils.RespFail(r)
}

func (*sys) UpdatePassword(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"userId":   "required",
		"password": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_user").Where("userId=?", data["userId"]).Update(g.Map{"password": data["password"]})
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) MenuPageList(r *ghttp.Request) {
	menus, err := g.Model("sys_menu").FindAll()
	utils.GetErrExit(r, err)

	utils.RespOkWithData(r, menus)
}

func (*sys) AddMenu(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"parentId": "required",
		//"path":     "required",
		"menuType": "required",
		"title":    "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_menu").Insert(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) DeleteMenu(r *ghttp.Request) {
	_, err := g.Model("sys_menu").Where("menuId=?", r.GetString("menuId")).Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) UpdateMenu(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"parentId": "required",
		//"path":     "required",
		"menuType": "required",
		"title":    "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_menu").Replace(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) RoleMenu(r *ghttp.Request) {
	roleId := r.GetString("roleId")

	// get role-menu all menu
	roleMenuList, err := g.Model("sys_role_menu").Where("roleId=?", roleId).FindAll()
	utils.GetErrExit(r, err)

	//var menuIds []string
	menuIds := utils.ResultStrSlice(roleMenuList, "menuIds")
	//for _, roleMenuObj := range roleMenuList {
	//	menuIds = append(menuIds, roleMenuObj["menuId"].String())
	//}

	menus, err := g.Model("sys_menu").FindAll()
	utils.GetErrExit(r, err)

	if menus.IsEmpty() {
		utils.RespOkWithData(r, []string{})
	}

	for _, menuObj := range menus {
		if utils.HasString(menuObj["menuId"].String(), menuIds) {
			menuObj["checked"] = gvar.New(true)
		}
	}

	utils.RespOkWithData(r, menus)
}

func (*sys) UpdateRoleMenu(r *ghttp.Request) {
	roleId := r.GetString("roleId")

	getJson, err := r.GetJson()
	utils.GetErrExit(r, err)

	menuIdInterfaces := getJson.Array()
	menuIds := utils.IS2SS(menuIdInterfaces)

	record, err := g.Model("sys_role_menu").Where("roleId=?", roleId).FindOne()
	utils.GetErrExit(r, err)

	if record.IsEmpty() {
		_, err = g.Model("sys_role_menu").Insert(g.Map{
			"roleId":  roleId,
			"menuIds": strings.Join(menuIds, ","),
		})
		utils.GetErrExit(r, err)
	} else {
		_, err = g.Model("sys_role_menu").Where("roleId=?", roleId).Update(g.Map{
			"menuIds": strings.Join(menuIds, ","),
		})
		utils.GetErrExit(r, err)
	}

	utils.RespOk(r)
}

func (*sys) DictList(r *ghttp.Request) {
	// get dictionaries
	dictionaries, err := g.Model("sys_dictionary").FindAll()
	utils.GetErrExit(r, err)

	if dictionaries.IsEmpty() {
		utils.RespOkWithData(r, []g.Map{})
	}

	utils.RespOkWithData(r, dictionaries)
}

func (*sys) AddDict(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"dictName": "required",
		"dictCode": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_dictionary").Insert(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) DeleteDict(r *ghttp.Request) {
	_, err := g.Model("sys_dictionary").Where("dictId=?", r.GetString("dictId")).Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) AddDataDict(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"dictName": "required",
		"dictCode": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_dictionary").Insert(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) DictPageList(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)

	data := r.GetMap()
	rules := map[string]string{
		"dictId": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	dictionaries, err := g.Model("sys_dictionary").Where("dictId=?", data["dictId"]).Offset((page - 1) * limit).Limit(limit).FindAll()
	utils.GetErrExit(r, err)

	if dictionaries.IsEmpty() {
		utils.RespOkWithData(r, []g.Map{})
	}

	total, err := g.Model("sys_dictionary").Count()
	utils.GetErrExit(r, err)

	utils.RespOkWithPageData(r, dictionaries, total)
}

func (*sys) UpdateDict(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"dictId":   "required",
		"dictName": "required",
		"dictCode": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_dictionary").Replace(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) DictDataPageList(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)

	data := r.GetMap()
	rules := map[string]string{
		"dictId": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	dictionaries, err := g.Model("sys_dictionary_data").Where("dictId=?", data["dictId"]).Offset((page - 1) * limit).Limit(limit).FindAll()
	utils.GetErrExit(r, err)

	if dictionaries.IsEmpty() {
		utils.RespOkWithData(r, []g.Map{})
	}

	total, err := g.Model("sys_dictionary_data").Count()
	utils.GetErrExit(r, err)

	utils.RespOkWithPageData(r, dictionaries, total)
}

func (*sys) DictDataList(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"dictCode": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	dictionaryObj, err := g.Model("sys_dictionary").Where("dictCode=?", data["dictCode"]).FindOne()
	utils.GetErrExit(r, err)

	fmt.Println("####->", dictionaryObj)

	if !dictionaryObj.IsEmpty() {
		dictDataObj, err := g.Model("sys_dictionary_data").Where("dictId=?", dictionaryObj["dictId"]).FindAll()
		utils.GetErrExit(r, err)

		if dictDataObj.IsEmpty() {
			utils.RespOkWithData(r, []g.Map{})
		}

		utils.RespOkWithData(r, dictDataObj)
	}

	utils.RespOkWithData(r, nil)
}

func (*sys) AddDictData(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"dictId":       "required",
		"dictDataName": "required",
		"dictDataCode": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_dictionary_data").Insert(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) DeleteDictData(r *ghttp.Request) {
	_, err := g.Model("sys_dictionary_data").Where("dictDataId=?", r.GetString("dictDataId")).Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) UpdateDictData(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"dictId":       "required",
		"dictDataName": "required",
		"dictDataCode": "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_dictionary_data").Replace(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) BatchDeleteDictData(r *ghttp.Request) {
	getJson, err := r.GetJson()
	utils.GetErrExit(r, err)
	dictDataIdInters := getJson.Array()
	dictDataIds := utils.IS2SS(dictDataIdInters)

	_, err = g.Model("sys_dictionary_data").Where("dictDataId in (?)", dictDataIds).Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) OrganizationList(r *ghttp.Request) {
	organizations, err := g.Model("sys_organization").FindAll()
	utils.GetErrExit(r, err)

	if organizations.IsEmpty() {
		utils.RespOkWithData(r, []g.Map{})
	}

	utils.RespOkWithData(r, organizations)
}

func (*sys) AddOrganization(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"organizationName":     "required",
		"organizationFullName": "required",
		"organizationCode":     "required",
		"organizationType":     "required",
		"parentId":             "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_organization").Insert(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) UpdateOrganization(r *ghttp.Request) {
	data := r.GetMap()
	rules := map[string]string{
		"organizationId":       "required",
		"organizationName":     "required",
		"organizationFullName": "required",
		"organizationCode":     "required",
		"organizationType":     "required",
		"parentId":             "required",
	}

	checkErr := g.Validator().Rules(rules).CheckMap(data)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, err := g.Model("sys_organization").Replace(data)
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}

func (*sys) DeleteOrganization(r *ghttp.Request) {
	_, err := g.Model("sys_organization").Where("organizationId=?", r.GetString("organizationId")).Delete()
	utils.GetErrExit(r, err)

	utils.RespOk(r)
}
