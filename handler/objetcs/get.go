/**
*FileName: objetcs
*Create on 2018/11/20 上午12:00
*Create by mok
*/

package objetcs

import (
	"github.com/gin-gonic/gin"
	"dataServer/model/file"
	"fmt"
	"dataServer/conf"
)


//获取到文件信息和内容
func GetFile(c *gin.Context){
	username := c.Query("username")
	path := c.Query("path")
	filename := c.Query("filename")
	fileinfo,err := file.GetFileInfo(username,path,filename)
	if err != nil{
		c.JSON(500,gin.H{
			"message":"打开文件失败",
			"err":err.Error(),
		})
		return
	}
	c.JSON(200,gin.H{
		"filename":fileinfo.UserFileInfo.FileName,
		"cratetime":fileinfo.CreateTime,
		"content":fileinfo.Content,
		"username":fileinfo.Username,
		"sha":fileinfo.Sha,
	})
}


//获取文件列表下所有文件的信息，主要是通过文件路径来获取
func GetFiles(c *gin.Context){
	path := c.Query("path")
	username := c.Query("username")
	filesPath := fmt.Sprintf("%s/%s/%s",conf.ETCD_FILEPATH_PREFIX,username,path)
	if path == ""{
		c.JSON(400,gin.H{
			"message":"获取该目录下的文件列表失败",
			"error":"传入的目录参数错误",
		})
		return
	}
	//todo:通过etcd获取到该文件路径，然后打印出该文件路径下所有文件的信息
	fs,err := file.GetUserFileInfos(filesPath)
	if err != nil{
		c.JSON(500,gin.H{
			"message":"获取文件列表失败",
			"err":err.Error(),
		})
		return
	}
	c.JSON(200,gin.H{
		"message":"OK",
		"data":fs,
	})
}


//直接返回文件内容
func GetRealfile(c *gin.Context){
	fileUrl := c.Query("fileUrl")
	if fileUrl ==""{
		c.JSON(400,nil)
		return
	}
	content,err := file.GetFileinByOther(fileUrl)
	if err !=nil{
		c.JSON(400,nil)
		return
	}
	c.JSON(200,content)
}