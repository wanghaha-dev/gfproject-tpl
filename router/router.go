package router

import (
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"myproject/app/api"
)

func init() {
	s := g.Server()
	s.Use(func(r *ghttp.Request) {
		r.Response.CORSDefault()
		r.Middleware.Next()
	})

	//s.Use(utils.CacheMiddleware)

	// 无需登录的
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/hello", api.Hello)

	})

	// 需要登录才能访问的
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(JwtAuthMiddleware)

	})
}
