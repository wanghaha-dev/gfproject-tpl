package utils

import (
	"github.com/wanghaha-dev/gf/net/ghttp"
)

const (
	SUCCESS_STATUS = 0
	FAIL_STATUS    = -1

	SUCCESS_MSG = "success"
	FAIL_MSG    = "fail"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type pageResponse struct {
	Page      int         `json:"page"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Total     int         `json:"total"`
	TotalPage int         `json:"totalPage"`
	PageSize  int         `json:"pageSize"`
}

func RespOk(r *ghttp.Request) {
	resp := &response{
		Code:    SUCCESS_STATUS,
		Message: SUCCESS_MSG,
	}
	_ = r.Response.WriteJsonExit(resp)
}

func RespOkWithData(r *ghttp.Request, data interface{}) {
	resp := &response{
		Code:    SUCCESS_STATUS,
		Message: SUCCESS_MSG,
		Data:    data,
	}
	_ = r.Response.WriteJsonExit(resp)
}

func RespOkWithDataAndCache(r *ghttp.Request, data interface{}) {
	resp := &response{
		Code:    SUCCESS_STATUS,
		Message: SUCCESS_MSG,
		Data:    data,
	}

	// add cache
	AddCache(r, r.RequestURI, resp)
	_ = r.Response.WriteJsonExit(resp)
}

func RespOkWithPageData(r *ghttp.Request, data interface{}, total int) {
	page, pageSize := GetPageArgs(r)

	totalPage := total / pageSize
	if total%pageSize != 0 {
		totalPage += 1
	}

	resp := &pageResponse{
		Code:      SUCCESS_STATUS,
		Message:   SUCCESS_MSG,
		Data:      data,
		Total:     total,
		TotalPage: totalPage,
		PageSize:  pageSize,
		Page:      page,
	}
	_ = r.Response.WriteJsonExit(resp)
}

func RespOkWithPageDataAndCache(r *ghttp.Request, data interface{}, total int) {
	page, pageSize := GetPageArgs(r)

	totalPage := total / pageSize
	if total%pageSize != 0 {
		totalPage += 1
	}

	resp := &pageResponse{
		Code:      SUCCESS_STATUS,
		Message:   SUCCESS_MSG,
		Data:      data,
		Total:     total,
		TotalPage: totalPage,
		PageSize:  pageSize,
		Page:      page,
	}

	// add cache
	AddCache(r, r.RequestURI, resp)
	_ = r.Response.WriteJsonExit(resp)
}

func RespFail(r *ghttp.Request) {
	resp := &response{
		Code:    FAIL_STATUS,
		Message: FAIL_MSG,
	}
	_ = r.Response.WriteJsonExit(resp)
}

func RespFailWithMsg(r *ghttp.Request, msg string) {
	resp := &response{
		Code:    FAIL_STATUS,
		Message: msg,
	}
	_ = r.Response.WriteJsonExit(resp)
}

func RespFailWith401Msg(r *ghttp.Request, msg string) {
	resp := &response{
		Code:    401,
		Message: msg,
	}
	_ = r.Response.WriteJsonExit(resp)
}
