package api

import (
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"github.com/wanghaha-dev/gf/os/gtime"
	"myproject/utils"
)

var Hello = helloApi{}

type helloApi struct{}

// Index is a demonstration route handler for output "Hello World!".
func (*helloApi) Index(r *ghttp.Request) {
	utils.RespOkWithData(r, g.Map{
		"code": 0,
		"msg":  gtime.Datetime(),
	})
}
