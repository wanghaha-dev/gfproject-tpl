package router

import (
	"fmt"
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"myproject/utils"
	"myproject/utils/gjwt"
	"strings"
)

// AuthMiddleware 权限校验中间件
func AuthMiddleware(r *ghttp.Request) {
	fmt.Println("method:", r.Method, "r.Router.Uri:", r.Router.Uri)
	fmt.Println(fmt.Sprintf("%v>%v", r.Method, r.Router.Uri))

	tokenData, err := gjwt.ParseToken(r.GetHeader("Authorization"))
	if err != nil {
		_ = r.Response.WriteJsonExit(g.Map{"code": 401, "msg": err})
		return
	} else {
		if tokenData["username"] == "admin" {
			r.Middleware.Next()
			return
		}

		fmt.Println("tokenData:", tokenData)
		userObj, err := g.Model("sys_user").Where("userId=?", tokenData["userId"]).FindOne()
		utils.GetErrExit(r, err)

		rolesString := userObj["roles"].String()
		rolesSlice := strings.Split(rolesString, ",")
		fmt.Println("roleSlice:", rolesString)

		roleMenus, err := g.Model("sys_role_menu").Where("roleId in (?)", rolesSlice).FindAll()
		utils.GetErrExit(r, err)

		menuIds := utils.ResultStrSlice(roleMenus, "menuIds")

		menus, err := g.Model("sys_menu").Where("menuId in (?)", menuIds).FindAll()
		utils.GetErrExit(r, err)

		authorities := utils.ResultStrSlice(menus, "authority")

		method := r.Method
		router := r.Router.Uri
		curAuth := fmt.Sprintf("%v>%v", method, router)

		if utils.HasString(curAuth, authorities) {
			r.Middleware.Next()
		} else {
			//r.Middleware.Next()
			_ = r.Response.WriteJsonExit(g.Map{
				"code":    403,
				"message": "没有权限访问",
			})
		}
	}
}
