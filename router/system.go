package router

import (
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"myproject/app/api/sys"
)

func init() {
	s := g.Server()

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(AuthMiddleware)
		group.GET("/api/auth/user", sys.Sys.AuthUser)
	})

	s.Group("/api", func(group *ghttp.RouterGroup) {
		// verification code
		group.GET("/captcha", sys.Sys.Captcha)
		// login
		group.POST("/login", sys.Sys.Login)
	})

	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(AuthMiddleware)

		// system user
		group.GET("/system/user/page", sys.Sys.UserPageList)
		group.DELETE("/system/user/:userId", sys.Sys.DeleteUser)
		group.GET("/system/role", sys.Sys.RoleList)
		group.POST("/system/user", sys.Sys.AddUser)
		group.PUT("/system/user", sys.Sys.UpdateUser)
		group.GET("/system/user/existence", sys.Sys.Existence)
		group.PUT("/system/user/password", sys.Sys.UpdatePassword)
		group.DELETE("/system/user/batch", sys.Sys.BatchDeleteUser)

		// system role
		group.GET("/system/role/page", sys.Sys.RolePageList)
		group.PUT("/system/role", sys.Sys.UpdateRole)
		group.POST("/system/role", sys.Sys.AddRole)
		group.DELETE("/system/role/:roleId", sys.Sys.DeleteRole)
		group.DELETE("/system/role/batch", sys.Sys.BatchDeleteRole)

		// system menu
		group.GET("/system/menu", sys.Sys.MenuPageList)
		group.POST("/system/menu", sys.Sys.AddMenu)
		group.PUT("/system/menu", sys.Sys.UpdateMenu)
		group.DELETE("/system/menu/:menuId", sys.Sys.DeleteMenu)

		// system role-menu
		group.GET("/system/role-menu/:roleId", sys.Sys.RoleMenu)
		group.PUT("/system/role-menu/:roleId", sys.Sys.UpdateRoleMenu)

		// system dictionary
		group.GET("/system/dictionary", sys.Sys.DictList)
		group.POST("/system/dictionary", sys.Sys.AddDict)
		group.DELETE("/system/dictionary/:dictId", sys.Sys.DeleteDict)
		group.PUT("/system/dictionary", sys.Sys.UpdateDict)

		// system dictionary-data
		group.GET("/system/dictionary-data/page", sys.Sys.DictDataPageList)
		group.GET("/system/dictionary-data", sys.Sys.DictDataList)
		group.POST("/system/dictionary-data", sys.Sys.AddDictData)
		group.DELETE("/system/dictionary-data/:dictDataId", sys.Sys.DeleteDictData)
		group.PUT("/system/dictionary-data", sys.Sys.UpdateDictData)
		group.DELETE("/system/dictionary-data/batch", sys.Sys.BatchDeleteDictData)

		// system organization
		group.GET("/system/organization", sys.Sys.OrganizationList)
		group.POST("/system/organization", sys.Sys.AddOrganization)
		group.PUT("/system/organization", sys.Sys.UpdateOrganization)
		group.DELETE("/system/organization/:organizationId", sys.Sys.DeleteOrganization)
	})
}
