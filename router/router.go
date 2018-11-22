/**
*FileName: router
*Create on 2018/11/19 下午11:58
*Create by mok
*/

package router

import (
	"github.com/gin-gonic/gin"
	"dataServer/handler/objetcs"
)

func Load(router *gin.Engine){

	f :=router.Group("file")

	{
		f.GET("",objetcs.GetFile)
		f.POST("",objetcs.PostFile)
		//f.DELETE("/:objectsName")
	}

	fs := router.Group("files")
	{
		fs.GET("",objetcs.GetFiles)
	}

	router.GET("/realfile",objetcs.GetRealfile)
}