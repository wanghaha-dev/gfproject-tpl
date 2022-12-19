package artlib

import (
	"math/rand"
	"myproject/utils"
	"time"

	"github.com/wanghaha-dev/gf/container/gvar"
	"github.com/wanghaha-dev/gf/database/gdb"
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
	"github.com/wanghaha-dev/gf/util/guid"
)

var Alartist = new(alArtist)

type alArtist struct{}

// List 获取数据列表
func (*alArtist) List(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)

	total, err := g.DB("artlib").Model("al_artist").
		Where("publish_flag=?", "1").
		Count()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	result, err := g.DB("artlib").Model("al_artist").
		Where("publish_flag=?", "1").
		OrderDesc("create_date").
		Offset((page - 1) * limit).
		Limit(limit).
		FindAll()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	utils.RespOkWithPageDataAndCache(r, result, total)
}

// All 获取所有数据不分页
func (*alArtist) All(r *ghttp.Request) {
	result, err := g.DB("artlib").Model("al_artist").
		Where("publish_flag=?", "1").
		FindAll()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	utils.RespOkWithDataAndCache(r, result)
}

// Add 新增数据
func (*alArtist) Add(r *ghttp.Request) {
	params := r.GetMap()
	params["id"] = guid.S()

	rules, ruleErr := utils.GetAddFieldsRule("al_artist")
	if ruleErr != nil {
		g.Log().Error(ruleErr)
		utils.RespFailWithMsg(r, ruleErr.Error())
		return
	}

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, insertErr := g.DB("artlib").Model("al_artist").
		Insert(params)
	if insertErr != nil {
		g.Log().Error(insertErr)
		utils.RespFail(r)
	}

	utils.RespOk(r)
}

// DetailById 根据Id获取数据详细信息
func (*alArtist) DetailById(r *ghttp.Request) {
	record, err := g.DB("artlib").Model("al_artist").
		Where("publish_flag=?", "1").
		Where("id=?", r.GetString("id")).
		FindOne()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	firstWorkIds, err := g.DB("artlib").Model("al_work").
		Where("publish_flag=?", "1").
		Where(g.Map{
			"artist_ids like ?": "%" + r.GetString("id") + "%",
		}).Limit(10).FindAll()

	if len(firstWorkIds) == 0 {
		record["firstWork"] = gvar.New(nil)
	} else {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(len(firstWorkIds))
		firstWorkIdObj := firstWorkIds[n]

		if err != nil {
			g.Log().Error(err)
			utils.RespFail(r)
			return
		}
		firstWork, err := g.DB("artlib").Model("al_work_files").Where("work_id=?", firstWorkIdObj["id"].String()).FindOne()
		if err != nil {
			g.Log().Error(err)
			utils.RespFail(r)
			return
		}

		record["firstWork"] = gvar.New(firstWork)
	}

	utils.RespOkWithData(r, record)
}

// ArtistWorks 艺术家下的艺术作品
func (*alArtist) ArtistWorks(r *ghttp.Request) {
	artistId := r.GetString("id")
	page, limit := utils.GetPageArgs(r)

	count, err := g.Model("al_work").
		Where("artist_ids like ?", "%"+artistId+"%").Count()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	result, err := g.Model("al_work").
		Where("publish_flag=?", "1").
		Where("artist_ids like ?", "%"+artistId+"%").
		OrderDesc("create_date").
		Offset((page - 1) * limit).
		Limit(limit).
		FindAll()

	utils.RespOkWithPageData(r, result, count)
}

// DeleteById 根据Id删除数据
func (*alArtist) DeleteById(r *ghttp.Request) {
	_, err := g.DB("artlib").Model("al_artist").
		Where("id=?", r.GetString("id")).
		Delete()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	utils.RespOk(r)
}

// UpdateById 根据id修改数据
func (*alArtist) UpdateById(r *ghttp.Request) {
	// 新数据
	params := r.GetMap()

	rules, ruleErr := utils.GetUpdateFieldsRule("al_artist")
	if ruleErr != nil {
		g.Log().Error(ruleErr)
		utils.RespFailWithMsg(r, ruleErr.Error())
		return
	}

	checkErr := g.Validator().Rules(rules).CheckMap(params)
	if checkErr != nil {
		utils.RespFailWithMsg(r, checkErr.FirstString())
	}

	_, updateErr := g.DB("artlib").Model("al_artist").
		Where("publish_flag=?", "1").
		Where("id=?", r.GetString("id")).
		Update(params)
	if updateErr != nil {
		g.Log().Error(updateErr)
		utils.RespFail(r)
		return
	}

	utils.RespOk(r)
}

// Index 首页
func (*alArtist) Index(r *ghttp.Request) {
	// todo
	utils.RespOkWithData(r, "hello alArtist")
}

// TypesCount 获取分类下的艺术品数量
func (*alArtist) TypesCount(r *ghttp.Request) {
	artistId := r.GetString("id")

	// 拿到所有分类
	typeData, err := g.Model("z_mini_type_codes").Fields("code, name").Where("parent_code=?", "02").FindAll()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	typeDataNotEmpty := make(gdb.Result, 0)

	// 遍历分类, 获取每个分类的作品数量
	// 并且获取每个分类的代表图片
	for _, item := range typeData {
		count, err := g.Model("al_work").
			Where("publish_flag=?", "1").
			Where("artist_ids like ?", "%"+artistId+"%").
			Where("subject_codes like ?", item["code"]).
			Count()
		if err != nil {
			g.Log().Error(err)
			utils.RespFail(r)
			return
		}

		item["count"] = gvar.New(count)

		// 如果分类下有作品的话, 就拿第一个作品
		if count != 0 {
			firstWork, err := g.Model("al_work").
				Where("publish_flag=?", "1").
				Where("artist_ids like ?", "%"+artistId+"%").
				Where("subject_codes like ?", item["code"]).FindOne()
			if err != nil {
				g.Log().Error(err)
				utils.RespFail(r)
				return
			}

			item["image"] = gvar.New(firstWork["imgs"])

			typeDataNotEmpty = append(typeDataNotEmpty, item)
		}
	}

	// 拿到总数
	total, err := g.Model("al_work").
		Where("publish_flag=?", "1").
		Where("artist_ids like ?", "%"+artistId+"%").
		Count()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	// 若总作品数为0则不需要加了
	if total != 0 {
		first, err := g.Model("al_work").
			Where("publish_flag=?", "1").
			Where("artist_ids like ?", "%"+artistId+"%").
			FindOne()
		if err != nil {
			g.Log().Error(err)
			utils.RespFail(r)
			return
		}

		allData := make(gdb.Record)
		allData["count"] = gvar.New(total)
		allData["name"] = gvar.New("全部")
		allData["image"] = gvar.New(first["imgs"])
		allData["code"] = gvar.New(-1)
		typeDataNotEmpty = append(typeDataNotEmpty, allData)
	}

	utils.RespOkWithData(r, typeDataNotEmpty)
}

// GetRelation 获取艺术家关联的艺术故事
func (*alArtist) GetRelation(r *ghttp.Request) {
	page, limit := utils.GetPageArgs(r)
	artistId := r.GetString("id")

	// 获取到所有的关联艺术故事id
	data, err := g.Model("al_appreciation_artist").Where("id=?", artistId).FindAll()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	// 根据故事id找到所有的艺术故事数据
	storyIds := utils.ResultStrSlice(data, "id")
	count, err := g.Model("al_appreciation").Where("id in ?", storyIds).Count()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	result, err := g.Model("al_appreciation").
		Offset((page-1)*limit).
		Limit(limit).
		Where("id in ?", storyIds).FindAll()
	if err != nil {
		g.Log().Error(err)
		utils.RespFail(r)
		return
	}

	utils.RespOkWithPageData(r, result, count)
}
