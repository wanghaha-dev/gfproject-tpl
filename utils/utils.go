package utils

import (
	"fmt"
	"github.com/wanghaha-dev/gf/database/gdb"
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"github.com/wanghaha-dev/gf/util/gconv"
	"strings"
)

func ResultStrSlice(result gdb.Result, key string) []string {
	var list []string
	for _, item := range result {
		menuIds := item[key].String()
		if menuIds != "" {
			list = append(list, strings.Split(menuIds, ",")...)
		}
	}
	return list
}

func GetErrExit(r *ghttp.Request, err error) {
	if err != nil {
		RespFail(r)
		fmt.Println("err:", err)
		g.Log().Error(err)
		return
	}
}

func GetErrWithMsgExit(r *ghttp.Request, err error, msg string) {
	if err != nil {
		RespFailWithMsg(r, msg)
		g.Log().Error(msg)
		return
	}
}

func HasString(val string, list []string) bool {
	has := false
	for _, item := range list {
		if item == val {
			has = true
		}
	}
	return has
}

func IS2SS(interfaces []interface{}) []string {
	var stringList []string
	for _, item := range interfaces {
		stringList = append(stringList, gconv.String(item))
	}
	return stringList
}
