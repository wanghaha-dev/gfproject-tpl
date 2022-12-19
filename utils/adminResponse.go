package utils

import (
	"github.com/wanghaha-dev/gf/database/gdb"
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
)

type AdminResponseFields struct {
	List      interface{} `json:"list"`
	Count     int         `json:"count"`
	TotalPage int         `json:"totalPage"`
	PageSize  int         `json:"pageSize"`
	Page      int         `json:"page"`
}

// AdminResponseFieldsEmpty 返回空
type AdminResponseFieldsEmpty struct {
	List      []g.Map `json:"list"`
	Count     int     `json:"count"`
	TotalPage int     `json:"totalPage"`
	PageSize  int     `json:"pageSize"`
	Page      int     `json:"page"`
}

type AdminResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    AdminResponseFields `json:"data"`
}

type AdminResponseEmpty struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    AdminResponseFieldsEmpty `json:"data"`
}

func AdminRespOkWithPageData(r *ghttp.Request, data interface{}, total int) {
	page, pageSize := GetPageArgs(r)

	totalPage := total / pageSize
	if total%pageSize != 0 {
		totalPage += 1
	}

	// 判断数据是否gdb.Result并且判断是否为空
	if data2, ok := data.(gdb.Result); ok {
		if data2.IsEmpty() {
			resp := &AdminResponseEmpty{
				Code:    SUCCESS_STATUS,
				Message: SUCCESS_MSG,
				Data: AdminResponseFieldsEmpty{
					List:      []g.Map{},
					Count:     total,
					TotalPage: totalPage,
					PageSize:  pageSize,
					Page:      page,
				},
			}
			_ = r.Response.WriteJsonExit(resp)
			return
		}
	}

	resp := &AdminResponse{
		Code:    SUCCESS_STATUS,
		Message: SUCCESS_MSG,
		Data: AdminResponseFields{
			List:      data,
			Count:     total,
			TotalPage: totalPage,
			PageSize:  pageSize,
			Page:      page,
		},
	}
	_ = r.Response.WriteJsonExit(resp)
}
