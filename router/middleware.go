package router

import (
	"github.com/wanghaha-dev/gf/net/ghttp"
	"myproject/utils"
	"myproject/utils/gjwt"
)

// JwtAuthMiddleware jwt身份验证中间件
func JwtAuthMiddleware(r *ghttp.Request) {
	token := r.GetHeader("token")
	_, err := gjwt.ParseToken(token)
	if err != nil {
		utils.RespUnauthorized(r)
		return
	}

	r.Middleware.Next()
}
