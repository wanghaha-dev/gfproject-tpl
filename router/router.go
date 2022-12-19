package router

import (
	"myproject/app/api"
	"myproject/app/api/artlib"
	"myproject/utils"

	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/net/ghttp"
)

func init() {
	s := g.Server()
	s.Use(func(r *ghttp.Request) {
		r.Response.CORSDefault()
		r.Middleware.Next()
	})

	s.Use(utils.CacheMiddleware)

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/hello", api.Hello)

		// alCountry
		group.GET("/alCountry/:id/detail", artlib.AlCountry.DetailById)
		group.DELETE("/alCountry/:id/delete", artlib.AlCountry.DeleteById)
		group.GET("/alCountry/list", artlib.AlCountry.List)
		group.GET("/alCountry/all", artlib.AlCountry.All)
		group.POST("/alCountry/add", artlib.AlCountry.Add)
		group.POST("/alCountry/:id/update", artlib.AlCountry.UpdateById)

		// alPeriod
		group.GET("/alPeriod/:id/detail", artlib.AlPeriod.DetailById)
		group.DELETE("/alPeriod/:id/delete", artlib.AlPeriod.DeleteById)
		group.GET("/alPeriod/list", artlib.AlPeriod.List)
		group.GET("/alPeriod/all", artlib.AlPeriod.All)
		group.POST("/alPeriod/add", artlib.AlPeriod.Add)
		group.POST("/alPeriod/:id/update", artlib.AlPeriod.UpdateById)

		// alStyle
		group.GET("/alStyle/:id/detail", artlib.AlStyle.DetailById)
		group.DELETE("/alStyle/:id/delete", artlib.AlStyle.DeleteById)
		group.GET("/alStyle/list", artlib.AlStyle.List)
		group.GET("/alStyle/all", artlib.AlStyle.All)
		group.POST("/alStyle/add", artlib.AlStyle.Add)
		group.POST("/alStyle/:id/update", artlib.AlStyle.UpdateById)

		// alArtist
		group.GET("/alArtist/:id/detail", artlib.Alartist.DetailById)
		group.DELETE("/alArtist/:id/delete", artlib.Alartist.DeleteById)
		group.GET("/alArtist/list", artlib.Alartist.List)
		group.GET("/alArtist/all", artlib.Alartist.All)
		group.POST("/alArtist/add", artlib.Alartist.Add)
		group.POST("/alArtist/:id/update", artlib.Alartist.UpdateById)
		group.GET("/alArtist/:id/works", artlib.Alartist.ArtistWorks)
		group.GET("/alArtist/:id/typesCount", artlib.Alartist.TypesCount)

		// alMuseum
		group.GET("/alMuseum/:id/detail", artlib.Almuseum.DetailById)
		group.DELETE("/alMuseum/:id/delete", artlib.Almuseum.DeleteById)
		group.GET("/alMuseum/list", artlib.Almuseum.List)
		group.GET("/alMuseum/all", artlib.Almuseum.All)
		group.POST("/alMuseum/add", artlib.Almuseum.Add)
		group.POST("/alMuseum/:id/update", artlib.Almuseum.UpdateById)

		// alStory
		group.GET("/alStory/:id/detail", artlib.Alstory.DetailById)
		group.DELETE("/alStory/:id/delete", artlib.Alstory.DeleteById)
		group.GET("/alStory/list", artlib.Alstory.List)
		group.GET("/alStory/all", artlib.Alstory.All)
		group.POST("/alStory/add", artlib.Alstory.Add)
		group.POST("/alStory/:id/update", artlib.Alstory.UpdateById)

		// alWork
		group.GET("/alWork/:id/detail", artlib.Alwork.DetailById)
		group.GET("/alWork/:id/download/list", artlib.Alwork.Download)
		group.DELETE("/alWork/:id/delete", artlib.Alwork.DeleteById)
		group.GET("/alWork/list", artlib.Alwork.List)
		group.GET("/alWork/all", artlib.Alwork.All)
		group.POST("/alWork/add", artlib.Alwork.Add)
		group.POST("/alWork/:id/update", artlib.Alwork.UpdateById)
		group.GET("/alWork/firstByCode", artlib.Alwork.FirstByCode)
		group.GET("/alWork/typesCount", artlib.Alwork.TypesCount)

		// alWorkFiles
		group.GET("/alWorkFiles/:id/detail", artlib.Alworkfiles.DetailById)
		group.DELETE("/alWorkFiles/:id/delete", artlib.Alworkfiles.DeleteById)
		group.GET("/alWorkFiles/list", artlib.Alworkfiles.List)
		group.GET("/alWorkFiles/all", artlib.Alworkfiles.All)
		group.POST("/alWorkFiles/add", artlib.Alworkfiles.Add)
		group.POST("/alWorkFiles/:id/update", artlib.Alworkfiles.UpdateById)
		group.GET("/alWorkFiles/:workId/dzi", artlib.Alworkfiles.DZI)

		// alTheme
		group.GET("/alTheme/:id/detail", artlib.Altheme.DetailById)
		group.DELETE("/alTheme/:id/delete", artlib.Altheme.DeleteById)
		group.GET("/alTheme/list", artlib.Altheme.List)
		group.GET("/alTheme/all", artlib.Altheme.All)
		group.POST("/alTheme/add", artlib.Altheme.Add)
		group.POST("/alTheme/:id/update", artlib.Altheme.UpdateById)

		// typeCodes
		group.GET("/typeCodes/:id/detail", artlib.Typecodes.DetailById)
		group.DELETE("/typeCodes/:id/delete", artlib.Typecodes.DeleteById)
		group.GET("/typeCodes/list", artlib.Typecodes.List)
		group.GET("/typeCodes/all", artlib.Typecodes.All)
		group.POST("/typeCodes/add", artlib.Typecodes.Add)
		group.POST("/typeCodes/:id/update", artlib.Typecodes.UpdateById)
	})
}
