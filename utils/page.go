package utils

import (
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
)

// GetPageArgs 获取分页参数
func GetPageArgs(r *ghttp.Request) (page, limit int) {
	page = r.GetInt("page", 1)
	limit = r.GetInt("limit", g.Cfg().GetInt("web.limit"))

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = g.Cfg().GetInt("web.limit")
	}

	return
}
