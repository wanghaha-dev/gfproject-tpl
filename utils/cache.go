package utils

import (
	"encoding/json"
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
)

// AddCache 添加缓存
func AddCache(r *ghttp.Request, key string, data interface{}) {
	cacheTime := g.Cfg().GetInt("web.cacheTime")

	// add cache
	_, err := g.Redis().
		Do("setex", key, cacheTime*60*60, data)
	GetErrExit(r, err)
}

// RespCache 从缓存里直接返回
func RespCache(r *ghttp.Request, key string) bool {
	// get cache
	exists, err := g.Redis().DoVar("EXISTS", key)
	GetErrExit(r, err)

	if exists.Int() == 1 {
		redisData, _ := g.Redis().DoVar("get", key)

		data := make([]map[string]interface{}, 0)
		err := json.Unmarshal(redisData.Bytes(), &data)
		GetErrExit(r, err)

		RespOkWithData(r, data)
		return true
	}

	return false
}

// GetCacheData 返回缓存数据
func GetCacheData(key string) (map[string]interface{}, error) {
	// get cache
	exists, err := g.Redis().DoVar("EXISTS", key)
	if err != nil {
		g.Log().Error(err)
		return nil, err
	}

	if exists.Int() == 1 {
		redisData, _ := g.Redis().DoVar("get", key)

		data := make(map[string]interface{}, 0)
		err := json.Unmarshal(redisData.Bytes(), &data)
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	return nil, nil
}

// CacheMiddleware 缓存中间件
func CacheMiddleware(r *ghttp.Request) {
	// get cache ==================================================
	cacheCode := r.Request.RequestURI
	data, err := GetCacheData(cacheCode)
	GetErrExit(r, err)

	if data != nil {
		_ = r.Response.WriteJsonExit(data)
		return
	}
	// end cache ==================================================

	r.Middleware.Next()
}
